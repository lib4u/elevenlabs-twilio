# 📞 Go ElevenLabs-Twilio Integration

This project is a voice assistant implementation in Go that integrates **Twilio Voice API** with **ElevenLabs**. It initiates an outbound call using Twilio, connects to **Media Streams**, and forwards the audio stream to **ElevenLabs** via WebSocket. It also receives real-time speech transcription.

## 🔧 Current Features

- 📲 Initiates **outbound phone calls** via Twilio REST API  
- 🔁 On connection:
  - Establishes **WebSocket connection with Twilio Media Streams**
  - Connects to **ElevenLabs via WebSocket**
  - **Forwards audio stream** from Twilio to ElevenLabs
- 🗣 Receives **real-time speech transcription**

## 📦 Technologies Used

- **Go**
- **Twilio Voice API** (outbound calls, Media Streams)
- **ElevenLabs API** (TTS/STT via WebSocket)
- **WebSocket** (bidirectional audio streaming)
