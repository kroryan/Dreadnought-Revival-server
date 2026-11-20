package matchmaker

// offlineTeamFor applies the host's optional team split to each human in a
// hosted match. DN_OFFLINE_TEAMS is "n:m", where n and m are the target human
// counts on Teams 1 and 2. Co-op modes keep everyone on Team 1. When the
// setting is absent, the matchmaker's normal group placement is preserved.
// The legacy values are retained for compatibility with existing launchers.
import (
	"os"
	"strconv"
	"strings"
)

func offlineTeamsMode() string {
	return strings.ToLower(strings.TrimSpace(os.Getenv("DN_OFFLINE_TEAMS")))
}

func offlineTeamFor(pvp bool, groupTeam int, perTeam map[int]int) int {
	if !pvp {
		return 1
	}
	mode := offlineTeamsMode()
	target1, target2, configured := 0, 0, false
	if parts := strings.Split(mode, ":"); len(parts) == 2 {
		n1, err1 := strconv.Atoi(parts[0])
		n2, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil && n1 >= 0 && n1 <= 5 && n2 >= 0 && n2 <= 5 && n1+n2 > 0 {
			target1, target2, configured = n1, n2, true
		}
	}
	if !configured {
		switch mode {
		case "together":
			target1, target2, configured = 5, 0, true
		case "versus":
			target1, target2, configured = 1, 1, true
		case "auto":
			target1, target2, configured = 5, 5, true
		}
	}
	if !configured {
		return groupTeam
	}
	if perTeam[1] < target1 {
		return 1
	}
	if perTeam[2] < target2 {
		return 2
	}
	if perTeam[1] <= perTeam[2] {
		return 1
	}
	return 2
}
