# Roadmap

This document represents draft development plans.


1. [x] [Create function that creates player](#create-player)
2. [x] [Create function that creates location](#create-location)
3. [x] [Create function that moves player to location](#move-player-to-location)
4. [x] [Create function that creates NPC](#create-npc)
5. [x] [Create function that starts conversation between player and NPC](#start-conversation)
6. [ ] [Respond to conversation](#respond-to-conversation)


## Feature requirements

### Create player

Create player with passed name. You must be able to get that name from player any time needed.

### Create location

Create location with passed name. You must be able to get that name from location any time needed.

### Move player to location

Move player to location.  
There must be another function to get player's current location.

### Create NPC

Create NPC with passed name. You must be able to get that name from NPC any time needed.

### Start conversation

Conversation is just a text that is returned on players attempt to talk with NPC.  
Extend NPC creation with additional value which will be returned every time player attempts to start conversation with him.  

### Respond to conversation
Player should be able to respond to conversation if it has such options.  
Option should lead to another conversation, making possible chaining deep conversation (dialogue). 