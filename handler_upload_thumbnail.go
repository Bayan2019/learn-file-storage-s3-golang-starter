package main

import (
	"fmt"
	"io"
	"net/http"

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
	mediaType := header.Header.Get("Content-Type")
	if mediaType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for thumbnail", nil)
		return
	}

	// Read all the image data into a byte slice using io.ReadAll
	data, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading file", err)
		return
	}

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
	videoThumbnails[videoID] = thumbnail{
		data:      data,
		mediaType: mediaType,
	}

	url := fmt.Sprintf("http://localhost:%s/api/thumbnails/%s", cfg.port, videoID)
	video.ThumbnailURL = &url

	// Update the database so that the existing video record has a new thumbnail URL
	// by using the cfg.db.UpdateVideo function.
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	// Respond with updated JSON of the video's metadata.
	// Use the provided respondwithJSON function and pass it the updated database.Video struct to marshal.
	respondWithJSON(w, http.StatusOK, video)

	// respondWithJSON(w, http.StatusOK, struct{}{})
}
