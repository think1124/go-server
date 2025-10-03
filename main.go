// 파일 이름: main.go

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// --- 여기부터는 실제 DB를 쓴다고 가정하고 만든 임시 데이터와 함수들입니다. ---
// --- 나중에 실제 DB로 이 부분만 교체하면 됩니다. ---

// 임시 유저 데이터 구조체
type TempUser struct {
	ZepetoID string
	Username string
	Count    int
}

// 임시 데이터베이스 역할 (실제로는 DB에 저장된 테이블)
var tempDatabase = []TempUser{
	{ZepetoID: "zepeto_god", Username: "ZEPETO_GOD", Count: 99999},
	{ZepetoID: "world_master", Username: "WorldMaster", Count: 88888},
	{ZepetoID: "dev_joy", Username: "개발하는조이", Count: 76543},
	{ZepetoID: "coding_fun", Username: "코딩조아", Count: 65432},
	{ZepetoID: "my_zepeto_id", Username: "MyNickName", Count: 54321}, // 테스트용 내 ID
	{ZepetoID: "user_a", Username: "유저A", Count: 43210},
	{ZepetoID: "user_b", Username: "유저B", Count: 32109},
	{ZepetoID: "user_c", Username: "유저C", Count: 21098},
	{ZepetoID: "user_d", Username: "유저D", Count: 10987},
	{ZepetoID: "user_e", Username: "유저E", Count: 9876},
}

// 임시 DB에서 랭킹을 조회하는 함수 (실제 DB 함수를 흉내 낸 것)
func getRankingsFromTempDB() []TempUser {
	// 글자 수(Count) 기준으로 내림차순 정렬
	sort.Slice(tempDatabase, func(i, j int) bool {
		return tempDatabase[i].Count > tempDatabase[j].Count
	})
	return tempDatabase
}

// ------------------- 임시 데이터베이스 코드 끝 -------------------


// ZEPETO 클라이언트에 보낼 랭킹 데이터 구조체 (필드명 최소화)
type RankEntry struct {
	Rank  int    `json:"r"` // Rank
	User  string `json:"u"` // Username
	Count int    `json:"c"` // Count
}

// 최종 응답 구조체 (전체 랭킹 + 내 순위)
type RankingResponse struct {
	TopRankings []RankEntry `json:"top"`
	MyRank      *RankEntry  `json:"myRank,omitempty"`
}

// ZEPETO 랭킹 API 핸들러
func zepetoRankingHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")

	// 1. (임시) DB에서 랭킹 데이터를 가져옵니다.
	allRankings := getRankingsFromTempDB()

	// 2. 전체 랭킹 목록(Top 50)과 내 랭킹 정보를 찾습니다.
	var topRankings []RankEntry
	var myRank *RankEntry = nil
	
	for i, user := range allRankings {
		rank := i + 1
		// 상위 50위까지만 topRankings에 추가
		if rank <= 50 {
			topRankings = append(topRankings, RankEntry{
				Rank:  rank,
				User:  user.Username,
				Count: user.Count,
			})
		}
		// 내 랭킹 찾기 (대소문자 무시)
		if userID != "" && strings.EqualFold(user.ZepetoID, userID) {
			myRank = &RankEntry{
				Rank:  rank,
				User:  user.Username,
				Count: user.Count,
			}
		}
	}

	// 3. 최종 응답 데이터를 만듭니다.
	response := RankingResponse{
		TopRankings: topRankings,
		MyRank:      myRank,
	}

	// 4. JSON으로 변환하여 응답합니다.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := chi.NewRouter()

	// Gzip 압축 미들웨어를 적용합니다.
	r.Use(middleware.Compress(5, "application/json")) 

	// ZEPETO 랭킹 API 주소를 등록합니다.
	r.Get("/zepeto/rankings", zepetoRankingHandler)

	// 서버 시작
	port := "8080"
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}