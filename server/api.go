package server

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/tuommii/jumbo/model"
)

func (s *Server) apiHome(w http.ResponseWriter, r *http.Request) {
	// ignore error, new session is always returned
	session, _ := s.cookies.Get(r, "mysession")

	flashes := session.Flashes()
	session.Save(r, w)

	var msg string

	if len(flashes) > 0 {
		flash := flashes[0].(string)
		msg = flash
	}

	players, err := s.db.GetPlayers()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	games, err := s.db.GetGames()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Players []model.Player
		Games   []model.Game
		Message string
	}{
		players,
		games,
		msg,
	}
	s.templates["home.html"].ExecuteTemplate(w, "base", response)
}

func (s *Server) apiCreateMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, _ := s.cookies.Get(r, "mysession")

	gameName := r.FormValue("gameName")
	winner := r.FormValue("winner")
	loser := r.FormValue("loser")
	comment := r.FormValue("comment")
	isTieStr := r.FormValue("isTie")

	if winner == "" || loser == "" {
		http.Error(w, "winner and loser required", http.StatusInternalServerError)
		return
	}

	if strings.EqualFold(winner, loser) {
		http.Error(w, "winner and loser cant be same player", http.StatusInternalServerError)
		return
	}

	// SQL might cry for empty strings
	if comment == "" {
		comment = "EMPTY"
	}

	var isTie bool
	if isTieStr == "tie" {
		isTie = true
	}

	match := model.Match{
		GameName: gameName,
		Winner:   winner,
		Loser:    loser,
		Comment:  comment,
		IsTie:    isTie,
	}

	_, err := s.db.CreateMatch(match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.AddFlash("Added new game: " + gameName + " | " + winner + " - " + loser)
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) apiCreatePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	playerName := r.FormValue("playerName")

	_, err := s.db.CreatePlayer(playerName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) apiCreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameName := r.FormValue("gameName")

	_, err := s.db.CreateGame(gameName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) apiDeleteMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	log.Println("ID:", id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	num, err := s.db.DeleteMatch(id)
	log.Println("NUM", num)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) apiDeleteGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameName := r.FormValue("gameName")

	_, err := s.db.DeleteGame(gameName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) apiDeletePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	playerName := r.FormValue("playerName")

	_, err := s.db.DeletePlayer(playerName)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) apiSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameName := r.FormValue("gameName")
	p1 := r.FormValue("player1")
	p2 := r.FormValue("player2")

	ld, err := strconv.Atoi(r.FormValue("limitDays"))
	if err != nil {
		ld = 0
	}

	lg, err := strconv.Atoi(r.FormValue("limitGames"))
	if err != nil {
		lg = 0
	}

	// Create filter based on form values
	f := model.Filter{
		GameName:   gameName,
		Player1:    p1,
		Player2:    p2,
		LimitDays:  ld,
		LimitGames: lg,
	}

	matches, err := s.db.GetMatches(f)
	if err != nil {
		log.Println(err)
		http.Error(w, "Search error", http.StatusInternalServerError)
		return
	}

	// Calculate stats form matches
	stats, err := model.StatsFromMatches(matches)
	if err != nil {
		log.Println(err)
		http.Error(w, "Stats error", http.StatusInternalServerError)
		return
	}

	sort.Sort(stats)

	// Anonyme struct
	data := struct {
		Stats    model.SortedStats
		Matches  []model.Match
		GameName string
	}{
		stats,
		matches,
		f.GameName,
	}
	// json.NewEncoder(w).Encode(data)
	// http.Redirect(w, r, "/api/results", http.StatusSeeOther)
	s.templates["results.html"].ExecuteTemplate(w, "base", data)
}

// Auth middleware
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Prompt credentials in browser
		w.Header().Set("WWW-Authenticate", `Basic realm="Jumbo - Track Stats"`)

		user, pswd, _ := r.BasicAuth()
		if user != username || pswd != password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// Favicon server favicon
func (s *Server) favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.png")
}
