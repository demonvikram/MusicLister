package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	ID         string     `json:"id"`
	SecretCode string     `json:"secret_code"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Playlists  []Playlist `json:"playlists"`
}

type Playlist struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Songs []Song `json:"songs"`
}

type Song struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Composer string `json:"composer"`
	MusicURL string `json:"music_url"`
}

var users []User

func main() {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/viewProfile", viewProfileHandler).Methods("GET")
	r.HandleFunc("/getAllSongsOfPlaylist", getAllSongsOfPlaylistHandler).Methods("GET")
	r.HandleFunc("/createPlaylist", createPlaylistHandler).Methods("POST")
	r.HandleFunc("/addSongToPlaylist", addSongToPlaylistHandler).Methods("POST")
	r.HandleFunc("/deleteSongFromPlaylist", deleteSongFromPlaylistHandler).Methods("DELETE")
	r.HandleFunc("/deletePlaylist", deletePlaylistHandler).Methods("DELETE")
	r.HandleFunc("/getSongDetail", getSongDetailHandler).Methods("GET")

	// Start the HTTP server
	r.HandleFunc("/", serveHome).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>hello world this is a music APP</h1>"))
}

// login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	secretCode := r.FormValue("secret_code")
	user, found := findUserBySecretCode(secretCode)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, user)
}

// register handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")

	if name == "" || email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	secretCode := generateUniqueSecretCode()
	user := User{
		ID:         generateUniqueID(),
		SecretCode: secretCode,
		Name:       name,
		Email:      email,
		Playlists:  []Playlist{},
	}
	users = append(users, user)
	jsonResponse(w, user)
}

// func viewProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Print("hello users")
// 	jsonResponse(w, users)
// }

//view userProfile
func viewProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	user, found := findUserByID(userID)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, user)
}


//get all songs of particular playlist
func getAllSongsOfPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	playlistID := r.FormValue("playlist_id")
	playlist, found := findPlaylistByID(playlistID)
	if !found {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, playlist.Songs)
}


//creating playlist
func createPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	playlistName := r.FormValue("name")

	if playlistName == "" {
		http.Error(w, "Playlist name is required", http.StatusBadRequest)
		return
	}

	user, found := findUserByID(userID)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// playlist := Playlist{
	// 	ID:    generateUniqueID(),
	// 	Name:  playlistName,
	// 	Songs: []Song{},
	// }

	// user.Playlists = append(user.Playlists, playlist)
	// jsonResponse(w, playlist)
	newPlaylist := Playlist{
		ID:    generateUniqueID(),
		Name:  playlistName,
		Songs: []Song{},
	}

	// Add the new playlist to the user's playlists
	user.Playlists = append(user.Playlists, newPlaylist)
	ans := deleteUserByID(user.ID)
	if ans {
		users = append(users, user)
	}
	jsonResponse(w, user)
}

//song add kr do playlist m
func addSongToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	playlistID := r.FormValue("playlist_id")
	songName := r.FormValue("song_name")
	composer := r.FormValue("composer")
	musicURL := r.FormValue("music_url")

	if songName == "" || composer == "" || musicURL == "" {
		http.Error(w, "Song name, composer, and music URL are required", http.StatusBadRequest)
		return
	}

	playlist, found := findPlaylistByID(playlistID)
	if !found {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	song := Song{
		ID:       generateUniqueID(),
		Name:     songName,
		Composer: composer,
		MusicURL: musicURL,
	}

	user, found := findUserByID(userId)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//user.Playlists = nil
	playlist.Songs = append(playlist.Songs, song)
	// user.Playlists := append(user.Playlists, playlist)

	for i, playlist := range user.Playlists {
		if playlist.ID == playlistID {
			user.Playlists = append(user.Playlists[:i], user.Playlists[i+1:]...)
			//deleteUserByID(user.ID)
			//users = append(users, user)
			//jsonResponse(w, map[string]string{"message": "Playlist deleted successfully"})
			//return
		}
	}
	
	user.Playlists = append(user.Playlists, playlist)

	ans := deleteUserByID(userId)
	if ans {
		users = append(users, user)
	}

	jsonResponse(w, user)
}

func deleteSongFromPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	playlistID := r.FormValue("playlist_id")
	songID := r.FormValue("song_id")

	user, found := findUserByID(userID)
	if !found {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	playlist, found := findPlaylistByID(playlistID)
	if !found {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	for i, song := range playlist.Songs {
		if song.ID == songID {
			playlist.Songs = append(playlist.Songs[:i], playlist.Songs[i+1:]...)
			for idx, curr_playlist := range user.Playlists {
				if curr_playlist.ID == playlistID {
					user.Playlists = append(user.Playlists[:idx], user.Playlists[idx+1:]...)
					user.Playlists=append(user.Playlists,playlist)
				}
			}
			deleteUserByID(userID)
			users = append(users, user)
			jsonResponse(w, map[string]string{"message": "Song deleted successfully"})
			return
		}
	}

	http.Error(w, "Song not found in the playlist", http.StatusNotFound)
}

func deletePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	playlistID := r.FormValue("playlist_id")

	user, found := findUserByID(userID)
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for i, playlist := range user.Playlists {
		if playlist.ID == playlistID {
			user.Playlists = append(user.Playlists[:i], user.Playlists[i+1:]...)
			deleteUserByID(user.ID)
			users = append(users, user)
			jsonResponse(w, map[string]string{"message": "Playlist deleted successfully"})
			return
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)
}

func getSongDetailHandler(w http.ResponseWriter, r *http.Request) {
	songID := r.FormValue("song_id")

	for _, user := range users {
		for _, playlist := range user.Playlists {
			for _, song := range playlist.Songs {
				if song.ID == songID {
					jsonResponse(w, song)
					return
				}
			}
		}
	}

	http.Error(w, "Song not found", http.StatusNotFound)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
	}
}

// Helper functions
func generateUniqueID() string {
	id := uuid.New()
	return id.String()
}

func generateUniqueSecretCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}

func findUserBySecretCode(secretCode string) (User, bool) {
	for _, user := range users {
		if user.SecretCode == secretCode {
			return user, true
		}
	}
	return User{}, false
}

func findUserByID(userID string) (User, bool) {
	for _, user := range users {
		if user.ID == userID {
			return user, true
		}
	}
	return User{}, false
}

// Helper function to find a playlist by ID
func findPlaylistByID(playlistID string) (Playlist, bool) {
	for _, user := range users {
		for _, playlist := range user.Playlists {
			if playlist.ID == playlistID {
				return playlist, true
			}
		}
	}
	return Playlist{}, false
}

func deleteUserByID(userID string) bool {
	for i, user := range users {
		if user.ID == userID {
			users = append(users[:i], users[i+1:]...)
			return true
		}
	}
	return false
}
