package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	// validate the request
	videoIDString := r.PathValue("videoID")
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

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// Set a const maxMemory to 10MB.
	// I just bit-shifted the number 10 to the left 20 times
	// to get an int that stores the proper number of bytes.
	const maxMemory = 10 << 20 // 10 MB
	// Use (http.Request).ParseMultipartForm with the maxMemory const as an argument
	r.ParseMultipartForm(maxMemory)

	// Use r.FormFile to get the file data. The key the web browser is using is called "thumbnail"
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	// Get the media type from the file's Content-Type header
	// Use the mime.ParseMediaType function to get the media type from the Content-Type header
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type", err)
		return
	}
	// If the media type isn't either image/jpeg or image/png,
	// respond with an error (respondWithError helper)
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}

	// Use the Content-Type header to determine the file extension
	// Use the videoID to create a unique file path.
	assetPath := getAssetPath(videoID, mediaType)
	assetDiskPath := cfg.getAssetDiskPath(assetPath)

	// Use os.Create to create the new file
	dst, err := os.Create(assetDiskPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file on server", err)
		return
	}
	defer dst.Close()
	// Copy the contents from the multipart.File to the new file on disk using io.Copy
	if _, err = io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving file", err)
		return
	}

	// Read all the image data into a byte slice using io.ReadAll
	// data, err := io.ReadAll(file)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Error reading file", err)
	// 	return
	// }

	// Get the video's metadata from the SQLite database.
	// The apiConfig's db has a GetVideo method you can use
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find video", err)
		return
	}
	if video.UserID != userID {
		// If the authenticated user is not the video owner,
		// return a http.StatusUnauthorized response
		respondWithError(w, http.StatusUnauthorized, "Not authorized to update this video", nil)
		return
	}

	// Create a new thumbnail struct with the image data and media type
	// Add the thumbnail to the global map, using the video's ID as the key
	// videoThumbnails[videoID] = thumbnail{
	// 	data:      data,
	// 	mediaType: mediaType,
	// }

	// Use base64.StdEncoding.EncodeToString from the encoding/base64 package
	// to convert the image data to a base64 string.
	// base64Encoded := base64.StdEncoding.EncodeToString(data)

	// base64DataURL := fmt.Sprintf("data:%s;base64,%s", mediaType, base64Encoded)
	// Store the URL in the thumbnail_url column in the database.

	// Instead of encoding to base64, update the handler
	// to save the bytes to a file at the path /assets/<videoID>.<file_extension>
	url := cfg.getAssetURL(assetPath)
	video.ThumbnailURL = &url

	// Update the database so that the existing video record has a new thumbnail URL
	// by using the cfg.db.UpdateVideo function.
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		// delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	// Respond with updated JSON of the video's metadata.
	// Use the provided respondwithJSON function and pass it the updated database.Video struct to marshal.
	respondWithJSON(w, http.StatusOK, video)
}
