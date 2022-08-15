package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/tcolgate/mp3"
)

var joinedServers []string
var sounds []Sound

type Sound struct {
	message, filename string
	duration          time.Duration
}

func main() {

	// Load configuration file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection,", err)
	}

	// register commands
	sounds = append(sounds,
		Sound{"mo", "sounds/mo.mp3", time.Minute},
		Sound{"ma", "sounds/ma.mp3", time.Minute},
		Sound{"ginf", "sounds/ginf.mp3", time.Minute},
		Sound{"depression", "sounds/warum_bin_ich_so_fröhlich.mp3", time.Minute},
		Sound{"teams", "sounds/teams.mp3", time.Minute},
		Sound{"okay", "sounds/okay.mp3", time.Minute},
		Sound{"yeet", "sounds/yeet.mp3", time.Minute},
		Sound{"marcos", "sounds/marcos.mp3", time.Minute},
		Sound{"outlook", "sounds/outlook.mp3", time.Minute},
		Sound{"bonk", "sounds/bonk.mp3", time.Minute},
		Sound{"bruh", "sounds/bruh.mp3", time.Minute},
		Sound{"bann", "sounds/ban_den_weg.mp3", time.Minute},
		Sound{"jamoin", "sounds/ja_moin.mp3", time.Minute},
		Sound{"megalovania", "sounds/megalovania.mp3", time.Minute},
		Sound{"ough", "sounds/ough.mp3", time.Minute},
		Sound{"yooo", "sounds/yooooooooooo.mp3", time.Minute},
		Sound{"haha", "sounds/haha.mp3", time.Minute},
		Sound{"letsgo", "sounds/letsgo.mp3", time.Minute},
		Sound{"lugner", "sounds/minze.mp3", time.Minute},
		Sound{"minze", "sounds/minze.mp3", time.Minute},
		Sound{"electro", "sounds/electroboom.mp3", time.Minute},
		Sound{"boom", "sounds/full_bridge_rectifier_song.mp3", time.Minute},
		Sound{"rectify", "sounds/full_bridge_rectifier_song.mp3", time.Minute},
		Sound{"fichtl", "sounds/fichtl.mp3", time.Minute},
		Sound{"ara", "sounds/ara_ara.mp3", time.Minute},
		Sound{"amogus", "sounds/amogus.mp3", time.Minute},
		Sound{"donk", "sounds/donk.mp3", time.Minute},
		Sound{"brass", "sounds/brass.mp3", time.Minute},
		Sound{"gunga", "sounds/gunga.mp3", time.Minute},
		Sound{"wesgo", "sounds/wesgo.mp3", time.Minute},
		Sound{"boogie", "sounds/boogie.mp3", time.Minute},
		Sound{"laugh", "sounds/laugh.mp3", time.Minute},
		Sound{"splishsplash", "sounds/splishsplash.mp3", time.Minute},
		Sound{"onemorething", "sounds/onemorething.mp3", time.Minute},
		Sound{"jojo", "sounds/jojo.mp3", time.Minute},
		Sound{"maria", "sounds/maria.mp3", time.Minute},
		Sound{"gay", "sounds/imgay.mp3", time.Minute},
		Sound{"weeee", "sounds/weeee.mp3", time.Minute},
		Sound{"abadaba", "sounds/abadaba.mp3", time.Minute},
		Sound{"wah", "sounds/wah.mp3", time.Minute},
		Sound{"why", "sounds/why.mp3", time.Minute},
		Sound{"uno", "sounds/uno.mp3", time.Minute},
	)

	// figure out the duration of the sounds
	for i, sound := range sounds {
		t := 0.0

		r, err := os.Open(sound.filename)
		if err != nil {
			log.Fatal("Error opening sound file", err)
		}

		d := mp3.NewDecoder(r)
		var f mp3.Frame
		skipped := 0

		for {
			if err := d.Decode(&f, &skipped); err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal("Error decoding sound", err)
			}

			t = t + f.Duration().Seconds()
		}

		err = r.Close()
		if err != nil {
			log.Fatal("Error closing soundfile", err)
		}

		sounds[i].duration = time.Duration(t)*time.Second + time.Second
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close Discord session
	if err != nil {
		return
	}
	_ = dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.ToLower(m.Content) == "brelp" || strings.ToLower(m.Content) == "bruhelp" {
		var helpstring = sounds[0].message

		for i := 1; i < len(sounds); i++ {
			helpstring += fmt.Sprintf(", %s", sounds[i].message)
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, helpstring)
		return
	}

	// check for commands
	for _, sound := range sounds {
		if voiceMessageHandler(s, m, sound) {
			return
		}
	}
}

func voiceMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, sound Sound) bool {
	if strings.ToLower(m.Content) == sound.message {
		go playSound(s, m, sound.filename, sound)
		return true
	}
	return false
}

func playSound(s *discordgo.Session, m *discordgo.MessageCreate, filename string, sound Sound) {
	if contains(joinedServers, m.GuildID) {
		return
	}
	joinedServers = append(joinedServers, m.GuildID)

	// s.Guild() funktioniert hier nicht, weil die VoiceStates nur in "state-cached guilds" verfügbar sind,
	// deshalb s.State.Guild()
	st, _ := s.State.Guild(m.GuildID)

	vc := func() (vc *discordgo.Channel) {
		for _, state := range st.VoiceStates {
			if state.UserID == m.Author.ID {
				channel, _ := s.State.Channel(state.ChannelID)
				return channel
			}
		}
		return nil
	}()
	if vc == nil { // Do nothing if user is not in a voice channel
		joinedServers = remove(joinedServers, m.GuildID)
		return
	}

	dvc, err := s.ChannelVoiceJoin(vc.GuildID, vc.ID, false, true)
	if err != nil {
		log.Println("failed joining voice channel")
		joinedServers = remove(joinedServers, m.GuildID)
		return
	}

	go removeGuildAfterTimeout(m.GuildID, sound.duration, dvc)

	dgvoice.PlayAudioFile(dvc, filename, make(chan bool))
	err = dvc.Disconnect()
	if err != nil {
		log.Println("failed leaving voice channel")
	}

	joinedServers = remove(joinedServers, m.GuildID)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func remove(s []string, e string) []string {
	for i, v := range s {
		if v == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func removeGuildAfterTimeout(gid string, duration time.Duration, dvc *discordgo.VoiceConnection) {
	time.Sleep(duration)
	if contains(joinedServers, gid) {
		joinedServers = remove(joinedServers, gid)
	}
}
