package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	gospeckle "github.com/speckleworks/gospeckle/pkg"
	"github.com/spf13/cobra"
)

var clientID string

func init() {
	rootCmd.AddCommand(streamCmd)

	streamCmd.Flags().StringVarP(&streamID, "stream_id", "", "", "the ID of the stream to print events from")
	streamCmd.Flags().StringVarP(&clientID, "client_id", "", "", "the ID of the client reading/writing from/to the stream")
	streamCmd.MarkFlagRequired("stream_id")
}

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "print the data flowing through a stream",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Create new client if not specified
		if clientID == "" {
			clientRequest := gospeckle.APIClientRequest{
				StreamID: streamID,
			}
			newClient, _, err := speckleClient.APIClient.Create(ctx, clientRequest)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			clientID = newClient.ID
		}

		c, err := speckleClient.NewWebsocket(clientID, streamID)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		defer c.Close()
		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}

				if string(message) == "ping" {
					err := c.WriteMessage(websocket.TextMessage, []byte("alive"))
					if err != nil {
						log.Println(err)
					}
					return
				}

				log.Printf("recv: %s", message)
			}
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return

			// ===================================
			// An example of how to send a message
			// ===================================
			// case _ = <-ticker.C:

			// 	message := gospeckle.WebsocketMessage{
			// 		EventName: "broadcast",
			// 		StreamID:  streamID,
			// 		Payload:   "Look at all these chickens!",
			// 	}

			// 	err := c.WriteJSON(message)

			// 	if err != nil {
			// 		log.Println("write:", err)
			// 		return
			// 	}

			case <-interrupt:
				fmt.Println("Shutting off stream")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}

	},
}
