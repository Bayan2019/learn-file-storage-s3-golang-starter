package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

// Update the handlerUploadVideo handler code to store bucket and key as a comma delimited string in the video_url.
func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	// Set an upload limit of 1 GB (1 << 30 bytes)
	const uploadLimit = 1 << 30
	// using http.MaxBytesReader
	r.Body = http.MaxBytesReader(w, r.Body, uploadLimit)

	// Extract the videoID from the URL path parameters
	videoIDString := r.PathValue("videoID")
	// and parse it as a UUID
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	// Authenticate the user to get a userID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// Get the video metadata from the database,
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find video", err)
		return
	}
	// if the user is not the video owner, return a http.StatusUnauthorized response
	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to update this video", nil)
		return
	}

	// Parse the uploaded video file from the form data
	// Use (http.Request).FormFile with the key "video" to get a multipart.File in memory
	file, handler, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	// Remember to defer closing the file with (os.File).Close - we don't want any memory leakss
	defer file.Close()

	// Validate the uploaded file to ensure it's an MP4 video
	// Use mime.ParseMediaType and "video/mp4" as the MIME type
	mediaType, _, err := mime.ParseMediaType(handler.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type", err)
		return
	}
	if mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type, only MP4 is allowed", nil)
		return
	}

	// Save the uploaded file to a temporary file on disk.
	// Use os.CreateTemp to create a temporary file.
	tempFile, err := os.CreateTemp("", "tubely-upload.mp4")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create temp file", err)
		return
	}
	// defer remove the temp file with os.Remove
	defer os.Remove(tempFile.Name())
	// defer close the temp file (defer is LIFO, so it will close before the remove)
	defer tempFile.Close()

	// io.Copy the contents over from the wire to the temp file
	if _, err := io.Copy(tempFile, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not write file to disk", err)
		return
	}

	// Reset the tempFile's file pointer to the beginning with .Seek(0, io.SeekStart)
	// - this will allow us to read the file again from the beginning
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not reset file pointer", err)
		return
	}

	// to get the aspect ratio of the video file from the temporary file once it's saved to disk.
	directory := ""
	aspectRatio, err := getVideoAspectRatio(tempFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error determining aspect ratio", err)
		return
	}
	switch aspectRatio {
	case "16:9":
		directory = "landscape"
	case "9:16":
		directory = "portrait"
	default:
		directory = "other"
	}

	key := getAssetPath(mediaType)
	key = filepath.Join(directory, key)

	processedFilePath, err := processVideoForFastStart(tempFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing video", err)
		return
	}
	defer os.Remove(processedFilePath)

	processedFile, err := os.Open(processedFilePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not open processed file", err)
		return
	}
	defer processedFile.Close()

	// Put the object into S3 using PutObject.
	_, err = cfg.s3Client.PutObject(r.Context(), &s3.PutObjectInput{
		// The bucket name
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(key),
		Body:        tempFile,
		ContentType: aws.String(mediaType),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error uploading file to S3", err)
		return
	}

	// Update the VideoURL of the video record in the database with the S3 bucket and key.
	// to store bucket and key as a comma delimited string in the video_url
	// In handlerUploadVideo don't store the bucket and key as comma separated values in the video_url field.
	// Use your distribution's domain name, and then dynamically inject the S3 object's key.
	url := fmt.Sprintf("%s/%s", cfg.s3CfDistribution, key)
	video.VideoURL = &url
	// Store an actual URL again in the video_url column, but this time, use the cloudfront URL.
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}
	// Remove the dbVideoToSignedVideo method and all references to it
	// video, err = cfg.dbVideoToSignedVideo(video)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't generate presigned URL", err)
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, video)
}

// Create a function that takes a file path and returns the aspect ratio as a string
func getVideoAspectRatio(filePath string) (string, error) {
	// It should use exec.Command to run the same ffprobe
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)

	// Set the resulting exec.Cmd's Stdout field to a pointer to a new bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// .Run() the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %v", err)
	}

	// Unmarshal the stdout of the command from the buffer's .Bytes
	// into a JSON struct so that you can get
	var output struct {
		Streams []struct {
			// the width
			Width int `json:"width"`
			// and height fields.
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return "", fmt.Errorf("could not parse ffprobe output: %v", err)
	}

	if len(output.Streams) == 0 {
		return "", errors.New("no video streams found")
	}

	width := output.Streams[0].Width
	height := output.Streams[0].Height

	if width == 16*height/9 {
		return "16:9", nil
	} else if height == 16*width/9 {
		return "9:16", nil
	}
	return "other", nil
}

// Create a new function that takes a file path as input
// and creates and returns a new path to a file with "fast start" encoding
func processVideoForFastStart(filePath string) (string, error) {
	// Create a new string for the output file path
	// appended .processing to the input file
	newFilePath := fmt.Sprintf("%s.processing", filePath)
	// Create a new exec.Cmd using exec.Command
	exec := exec.Command("ffmpeg", "-i", filePath, "-movflags", "faststart",
		"-codec", "copy", "-f", "mp4", newFilePath)

	var stderr bytes.Buffer
	exec.Stderr = &stderr

	// Run the command
	if err := exec.Run(); err != nil {
		return "", fmt.Errorf("error processing video: %s, %v", stderr.String(), err)
	}

	fileInfo, err := os.Stat(newFilePath)
	if err != nil {
		return "", fmt.Errorf("could not stat processed file: %v", err)
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("processed file is empty")
	}

	// Return the output file path
	return newFilePath, nil
}

// Remove the dbVideoToSignedVideo method and all references to it
// func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
// 	if video.VideoURL == nil {
// 		return video, nil
// 	}
// 	// It should first split the video.VideoURL on the comma
// 	parts := strings.Split(*video.VideoURL, ",")
// 	if len(parts) < 2 {
// 		return video, nil
// 	}
// 	// to get the bucket
// 	bucket := parts[0]
// 	// and key
// 	key := parts[1]
// 	// use generatePresignedURL to get a presigned URL for the video
// 	// Remove generatePresignedURL function and all references to it

// 	// with the VideoURL field set to a presigned URL
// 	video.VideoURL = &presigned
// 	return video, nil
// }

// Remove generatePresignedURL function and all references to it

// func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
// 	// Use the SDK to create a s3.PresignClient with s3.NewPresignClient
// 	presignClient := s3.NewPresignClient(s3Client)
// 	// Use the client's .PresignGetObject() method with s3.WithPresignExpires as a functional option.
// 	presignedUrl, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(key),
// 	}, s3.WithPresignExpires(expireTime))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
// 	}
// 	// Return the .URL field of the v4.PresignedHTTPRequest created by .PresignGetObject()
// 	return presignedUrl.URL, nil
// }
