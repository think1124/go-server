// 파일 이름: main.go (HTML 방식 최종본 - embed 적용)

package main

import (
	"embed" // ★ 1. embed 패키지를 import 합니다.
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ★ 2. go:embed 지시자를 사용하여 templates 폴더의 내용을 변수에 담습니다.
// 이 코드는 컴파일 시점에 templates/ranking.html 파일을 templateFiles 변수로 읽어옵니다.
//go:embed templates/ranking.html
var templateFiles embed.FS

// 임시 유저 데이터 (이전과 동일)
type TempUser struct {
	ZepetoID string
	Username string
	Count    int
}

var tempDatabase = []TempUser{
	{ZepetoID: "zepeto_god", Username: "ZEPETO_GOD", Count: 99999},
	{ZepetoID: "world_master", Username: "WorldMaster", Count: 88888},
    // ... (이전과 동일한 임시 데이터)
}

// HTML 템플릿에 전달할 최종 랭킹 데이터 구조체
type PageRankEntry struct {
	Rank  int
	User  string
	Count int
}

// 랭킹 HTML 페이지를 보여주는 핸들러
func rankingHTMLHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 임시 DB에서 랭킹 데이터를 가져와 정렬합니다.
	sort.Slice(tempDatabase, func(i, j int) bool {
		return tempDatabase[i].Count > tempDatabase[j].Count
	})

	// 2. HTML 템플릿에 채워넣을 최종 데이터를 만듭니다.
	var pageRankings []PageRankEntry
	for i, user := range tempDatabase {
        if i >= 50 {
            break
        }
		pageRankings = append(pageRankings, PageRankEntry{
			Rank:  i + 1,
			User:  user.Username,
			Count: user.Count,
		})
	}

	// ★ 3. 파일을 직접 읽는 대신, embed로 주입된 변수에서 템플릿을 파싱합니다.
	tmpl, err := template.ParseFS(templateFiles, "templates/ranking.html")
	if err != nil {
		http.Error(w, "Could not parse template from embed", http.StatusInternalServerError)
		log.Printf("Template parse error: %v", err)
		return
	}

	// 4. HTML 파일에 랭킹 데이터를 채워서 사용자에게 전송합니다.
	err = tmpl.Execute(w, pageRankings)
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		log.Printf("Template execute error: %v", err)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/webranking", rankingHTMLHandler)

    // Render.com은 PORT 환경변수를 사용합니다. 없으면 8080을 사용하도록 설정.
	port := "8080" 
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}