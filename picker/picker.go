package picker

import (
	"fmt"
	"github.com/forquare/patreon_picker/utils"
	"gopkg.in/mxpv/patreon-go.v1"
	"strings"
	"time"
)

type Mention struct {
	Names       []string
	PublishDate string
	IsInPast    bool // Booleans default to False
}

func GetPatreonMentions(client *patreon.Client) []Mention {

	sundayDate := getLastSunday()
	isInPast := true

	patreons := fetchPatreons(client)
	numberOfWeeklyPatreons := getNumberOfWeeklyPatreons(len(patreons))

	numberOfWeeksToGenerate := 20
	mentions := make([]Mention, numberOfWeeksToGenerate)

	for i := 0; i < len(mentions); i++ {
		mentions[i] = makeMention(patreons, numberOfWeeklyPatreons, sundayDate, isInPast)
		isInPast = false
		sundayDate = sundayDate.AddDate(0, 0, 7)
	}

	return mentions
}

func fetchPatreons(client *patreon.Client) []string {
	// Get your campaign data
	campaignResponse, err := client.FetchCampaign()
	if err != nil {
		panic(err)
	}

	campaignId := campaignResponse.Data[0].ID
	patreonCount := campaignResponse.Data[0].Attributes.PatronCount

	cursor := ""
	page := 1
	patreons := make([]string, patreonCount+1)
	count := 0

	for {
		pledgesResponse, err := client.FetchPledges(campaignId,
			patreon.WithPageSize(25),
			patreon.WithCursor(cursor))

		if err != nil {
			panic(err)
		}

		// Get all the users in an easy-to-lookup way
		users := make(map[string]*patreon.User)
		for _, item := range pledgesResponse.Included.Items {
			u, ok := item.(*patreon.User)
			if !ok {
				continue
			}

			users[u.ID] = u

			if u.Attributes.FullName == utils.GetAuthenticatedUserName(client) {
				continue
			}

			patreons[count] = strings.TrimSpace(u.Attributes.FullName)
			count++
		}

		// Get the link to the next page of pledges
		nextLink := pledgesResponse.Links.Next
		if nextLink == "" {
			break
		}

		cursor = nextLink
		page++
	}

	var newPatreons []string

	for _, p := range patreons {
		if p != "" {
			newPatreons = append(newPatreons, p)
		}
	}

	patreons = newPatreons

	return patreons
}

func getNumberOfWeeklyPatreons(numberOfPatreons int) int {
	if numberOfPatreons == 13 {
		/*
			I don't like this, maybe this will help in the future?
			https://stackoverflow.com/questions/1748468/evenly-distribute-values-over-a-time-period
		*/
		return 2
	}

	weeksInYear := 52
	numberOfWeeklyPatreons := numberOfPatreons / weeksInYear

	if numberOfWeeklyPatreons < 2 {
		numberOfWeeklyPatreons = 2
	}

	for numberOfPatreons%numberOfWeeklyPatreons == 1 {
		numberOfWeeklyPatreons++
	}

	return numberOfWeeklyPatreons
}

func getLastSunday() time.Time {
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 1, 0, time.UTC)

	// Our first Sunday is 1 second in
	lastSunday := today.AddDate(0, 0, -(int(today.Weekday())))

	return lastSunday
}

func getWeekNumber(myDate time.Time) int {
	// Sunday 4th January 1970 00:00:01
	var firstSunday int64 = 259201

	/*
		Get the time in seconds since 1st January 1970
		Subtract the first four days so our base is a Sunday
		Reduce seconds to minutes, to hours, to days, then to weeks.

		This gives us the number of weeks since the first Sunday of 1970
	*/
	return int((myDate.Unix() - firstSunday) / 60 / 60 / 24 / 7)
}

func makeMention(allPatreons []string, numberOfWeeklyPatreons int, sundayDate time.Time, isInPast bool) Mention {
	var mention Mention

	daySuffix := "th"

	switch sundayDate.Day() {
	case 1, 21, 31:
		daySuffix = "st"
	case 2, 22, 32:
		daySuffix = "nd"
	case 3, 23, 33:
		daySuffix = "rd"
	default:
		daySuffix = "th"
	}

	mention.PublishDate = fmt.Sprintf("Sunday, %d%s %s %d", sundayDate.Day(), daySuffix, sundayDate.Month().String(), sundayDate.Year())
	mention.IsInPast = isInPast

	// Array to hold this week's Patreons
	var thisWeeksPatreons []string
	partitionSize := len(allPatreons) / numberOfWeeklyPatreons
	weekNumber := getWeekNumber(sundayDate)
	/*
		The really cool thing here is that weekNumber can (and will) be huge! Like several thousand.
		But by using modulus, we can make that represent a place between 0 and partitionSize that
		increments each week.
	*/
	cursor := weekNumber % partitionSize

	/*
		  Cursor                 Cursor+partitionSize
		      |                       |
		      v                       v
		+---+---+---+---+---+---+---+---+---+---+---+---+---+
		|   |   |   |   |   |   |   |   |   |   |   |   |   |
		| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10| 11| 12| 13|
		|   |   |   |   |   |   |   |   |   |   |   |   |   |
		+---+---+---+---+---+---+---+---+---+---+---+---+---+
		^ 0   1   2   3...  ^                   ^
		|                   |                   |
		|                   |                   |
		|                   | <partitionSize>   |
		+-------------------+-------------------+
		                    |
		                    |
		                 Partitions
	*/

	/*
		- Find the Patreon under the cursor and add it to []thisWeeksPatreons.
		- Increment the counter so that we add a new []thisWeeksPatreons
		- Increment the cursor by partitionSize so that next time we jump to the "same index"
		  in the next partition
	*/
	for cursor < len(allPatreons) {
		thisWeeksPatreons = append(thisWeeksPatreons, strings.TrimSpace(allPatreons[cursor]))
		cursor = cursor + partitionSize
	}

	mention.Names = thisWeeksPatreons

	return mention
}
