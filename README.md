# 🃏 S***head

## 📌 Overview
A turn-based multiplayer card game designed to be played with friends online. Built from the ground up using WebSockets for real-time interaction and custom game logic. This project is focused on implementing game state synchronization, connection handling, and a polished, cross-platform user interface.

## 🧪 Tech Stack
- **Backend:** Go (Gorilla WebSocket, Gorilla Mux)
- **Frontend:** React Native
- **Communication:** WebSockets
- **Optional Infra (Later):** Redis (game state cache), Docker, CI/CD

---

## 🎯 Goals
- ✅ Build a fully functional turn-based multiplayer game engine
- 🔗 Support real-time gameplay via WebSockets
- 📱 Cross-platform frontend for iOS, web, and iPad
- 🧪 Robust test coverage for game logic
- 🚀 Future goal: matchmaking and lobby system

---

## 🧱 Milestones

| Phase         | Description                             | Status         |
|---------------|-----------------------------------------|----------------|
| Design        | Game rules, UI layout, network schema   | 🟡 In Progress |
| Backend       | Game server, WebSocket events, state    | ⚪ Not Started  |
| Frontend      | Lobby, game board, turn handling        | ⚪ Not Started  |
| Testing       | Simulated gameplay tests, latency tests | ⚪ Not Started  |
| Multiplayer   | Connection management, reconnection     | ⚪ Not Started  |

---

## 🎨 Design

### 🧠 Game Logic Plan
- [ ] Define game flow (rounds, turns, win conditions)
- [ ] Enumerate card types and effects
- [ ] Map out game state transitions
- [ ] Build test cases for state machine

### 🖼️ UI Mockups
Designs will be uploaded as they are completed:

- `/assets/design/gameboard.png`
- `/assets/design/lobby.png`
- `/assets/design/turn-ui.png`

---

## 🚧 Development Roadmap

### 🎮 Game Design
- [x] Choose core rules and game structure
- [ ] Build a turn engine that validates moves
- [ ] Handle edge cases like disconnections, forfeits

### 🖧 Backend (Go)
- [x] Set up Gorilla WebSocket server
- [x] Implement join/create room functionality
- [x] Create match state and player state structs
- [x] Handle broadcast, receive, and update logic

### 📱 Frontend (React Native)
- [x] Build lobby/join screen
- [ ] Implement game screen with card layout
- [ ] Add user feedback (e.g., turn alerts, animations)

### 🧪 Testing
- [x] Unit tests for game logic
- [x] Integration tests for WebSocket flow
- [ ] Manual testing with multiple device emulators

### 🚀 Deployment
- [ ] Dockerize backend
- [ ] Explore Render or Fly.io for hosting
- [ ] Set up GitHub Actions for CI

---

## 🧠 Inspiration

This project stems from my girlfriend, who I would play this game with all the time. However, when we are away we don't have a way of playing, & the already existing online solutions aren't very user friendly or interactive feeling.

I also love strategic games & I believe that this would be a fun thing to build in order to explore Go's concurrency model & WebSocket architecture in a real-time setting.

---

## 📬 Contact

Interested in collaborating or playtesting?  
Reach out on [GitHub](https://github.com/tyreesamurai) or open an issue in the repo.
