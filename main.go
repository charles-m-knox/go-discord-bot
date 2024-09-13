package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Value is provided directly by a CLI flag.
var (
	token                    string
	presenceChangesPath      string
	msgLogPath               string
	eventsGuildIDs           string
	eventsPath               string
	eventsRefreshRateSeconds int
)

// Value is derived from a passed-in CLI flag.
var (
	presenceChangesBasePath string
	presenceChangesFileName string
	msgLogBasePath          string
	msgLogFileName          string
	channels                map[string]string
	guildEvents             map[string][]*discordgo.GuildScheduledEvent
)

// Runtime variable.
var (
	mu sync.Mutex
)

func init() {
	var chans string

	flag.StringVar(&presenceChangesPath, "p", path.Join("output", "presence.jsonl"), "Path to write presence changes to, such as 'presence.jsonl'.")
	flag.StringVar(&msgLogPath, "m", path.Join("output", "msgs.jsonl"), "Path to write message logs to, such as 'msgs.jsonl'.")
	flag.StringVar(&eventsPath, "ep", path.Join("output", "events"), "Base path to write message logs to. Events will be written to <guildId>.json under this path.")
	flag.StringVar(&token, "t", "", "Bot token")
	flag.StringVar(&chans, "c", "", "Semicolon-separated channel IDs and names that the bot will look for - for example: 1212159453190116000,test-guild#test-channel1;1212159453190116001,test-guild#test-channel2;")
	flag.StringVar(&eventsGuildIDs, "e", "", "Comma-separated guild IDs the bot retrieves events for - for example: 1212159453190116000,1212159453190116001")
	flag.IntVar(&eventsRefreshRateSeconds, "er", 3600, "Number of seconds to wait between scheduled event refreshes for every guild. Recommend using a value of 3600 or greater.")

	flag.Parse()

	splitChannels := strings.Split(chans, ";")

	channels = make(map[string]string)

	for _, c := range splitChannels {
		cc := strings.Split(c, ",")

		if len(cc) != 2 {
			continue
		}

		channels[cc[0]] = cc[1]
	}

	if len(channels) == 0 {
		log.Fatalf("no channels were provided, please use the the -c flag")
	}

	if presenceChangesPath == "" {
		log.Fatalf("presenceChangesPath must not be empty")
	}

	presenceChangesBasePath, presenceChangesFileName = path.Split(presenceChangesPath)

	if msgLogPath == "" {
		log.Fatalf("msgLogPath must not be empty")
	}

	msgLogBasePath, msgLogFileName = path.Split(msgLogPath)

	if eventsPath == "" {
		log.Fatalf("eventsPath must not be empty")
	}

	guildEvents = make(map[string][]*discordgo.GuildScheduledEvent)

	eg := strings.Split(eventsGuildIDs, ",")

	if len(eg) == 0 {
		log.Println("warning: no guild IDs were specified, no scheduled events will be processed")
	}

	for _, e := range eg {
		guildEvents[e] = []*discordgo.GuildScheduledEvent{}
	}
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(memberLeave)
	dg.AddHandler(memberJoin)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildPresences | discordgo.IntentsGuildMembers | discordgo.IntentsGuildScheduledEvents

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// start a goroutine to periodically retrieve the events for subscribed
	// servers
	go func() {
		sleepTime := time.Duration(eventsRefreshRateSeconds) * time.Second
		for {
			for guildID := range guildEvents {
				log.Printf("retrieving scheduled events for guild %v...", guildID)

				guildEvents[guildID], err = dg.GuildScheduledEvents(guildID, true)
				if err != nil {
					log.Printf("failed to get guild %v events: %v", guildID, err.Error())
				}

				log.Printf("writing %v scheduled events for guild %v...", len(guildEvents[guildID]), guildID)

				writeScheduledEvents(eventsPath, fmt.Sprintf("%v.json", guildID), guildEvents[guildID])
			}

			log.Printf("waiting %v before checking scheduled events again.", sleepTime)

			time.Sleep(sleepTime)
		}
	}()

	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
