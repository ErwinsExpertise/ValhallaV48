# Wedding System

## Implemented Flow

1. Engagement
- `9201000` crafts engagement ring boxes.
- A player uses the engagement box on another player on the same map.
- The target receives a marriage proposal popup.
- On acceptance:
  - both players are linked in `characters.partnerID`
  - a row is created in `marriages`
  - engagement ring items are granted to both players
- On decline:
  - the proposer receives the box back

2. Reservation
- Cathedral: `9201005`
- Chapel: `9201008`
- One partner must have the correct reservation ticket.
- Both partners must be on the same map/channel.
- Reservation grants 15 stacked invitation cards to each partner.

3. Guest Invitations
- Invitation cards are consumed through ring action mode `5`.
- The named guest receives the corresponding guest-entry invitation item.
- Opening the received invitation runs ring action mode `6` and shows the couple info packet.
- Wedding invite cards are stacked.

4. Ceremony Session
- Wedding sessions are tracked in memory per marriage.
- Session state includes:
  - venue
  - invitation item type
  - guest ticket type
  - entry map
  - altar map
  - invited guests
  - stage
  - blessing count
  - started/completed flags
- Stage progression is timed:
  - lobby
  - ceremony
  - blessings closed
  - party
  - cleanup/end

5. Wedding Start
- Cathedral officiator: `9201002`
- Chapel officiator: `9201012`
- Both partners must be present and still engaged.
- Engagement rings are consumed and wedding rings are granted.

6. Guest Entry
- Guests can enter through the assistant NPCs once a session is active.
- Entry goes to the lobby before ceremony and to the altar after ceremony start.

7. Divorce
- `9201004`
- Removes marriage linkage and clears spouse metadata.

8. Registry / Gifts
- Ring action mode `9` stores the couple's wishlist in the active reservation.
- Wedding action mode `6` sends gifts to the spouse registry.
- Wedding action mode `7` retrieves gifts from the player's registry queue.
- Wedding action mode `8` closes the registry interaction.

## NPCs / Maps Involved

### NPCs
- `9201000` engagement ring maker
- `9201002` Cathedral officiator
- `9201004` marriage info / divorce
- `9201005` Cathedral assistant
- `9201008` Chapel assistant
- `9201012` Chapel officiator
- `9201013` Cathedral info
- `9201014` registry / Onyx chest helper

### Maps
- `680000000` Amoria
- `680000100` Chapel waiting hall
- `680000110` Chapel altar
- `680000200` Cathedral waiting hall
- `680000210` Cathedral altar
- `680000300` wedding party / photo area
- `680000500` wedding exit

## NX Wedding Map Coverage

Ground truth was taken from local NX map placement data through Valhalla's NX loader.

### `680000100` Chapel Lounge
- Map name: `White Wedding Lounge`
- Street: `Amoria`
- NPC IDs present:
  - `9201008` Assistant Bonnie
  - `9201010` Assistant Travis
  - `9201014` Pila Present
- Scripts implemented or updated:
  - `9201008.js`
  - `9201010.js`
  - `9201014.js`
- Unresolved/missing NPCs:
  - none

### `680000110` Chapel Altar
- Map name: `White Wedding Altar`
- Street: `Amoria`
- NPC IDs present:
  - `9201009` Assistant Jackie
  - `9201010` Assistant Travis
  - `9201011` Pelvis Bebop
- Scripts implemented or updated:
  - `9201009.js`
  - `9201010.js`
  - `9201011.js`
- Unresolved/missing NPCs:
  - none

### `680000200` Cathedral Lounge
- Map name: `Saint Maple Lounge`
- Street: `Amoria`
- NPC IDs present:
  - `9201005` Assistant Nicole
  - `9201006` Assistant Debbie
  - `9201014` Pila Present
- Scripts implemented or updated:
  - `9201005.js`
  - `9201006.js`
  - `9201014.js`
- Unresolved/missing NPCs:
  - none

### `680000210` Cathedral Altar
- Map name: `Saint Maple Altar`
- Street: `Amoria`
- NPC IDs present:
  - `9201002` High Priest John
  - `9201006` Assistant Debbie
  - `9201007` Assistant Nancy
- Scripts implemented or updated:
  - `9201002.js`
  - `9201006.js`
  - `9201007.js`
- Unresolved/missing NPCs:
  - none

### `680000300` Cake Photo
- Map name: `Cherished Visage Photos`
- Street: `Amoria`
- NPC IDs present:
  - `9201010` Assistant Travis
  - `9201021` Robin The Huntress
- Scripts implemented or updated:
  - `9201010.js`
  - `9201021.js`
- Unresolved/missing NPCs:
  - none

### `680000400` Bonus Hunting Map
- Map name: `Untamed Hearts Hunting Ground`
- Street: `Amoria`
- NPC IDs present:
  - `9201021` Robin The Huntress
- Scripts implemented or updated:
  - `9201021.js`
- Unresolved/missing NPCs:
  - none

### `680000401` Bonus Gift Map
- Map name: `The Love Pinata`
- Street: `Amoria`
- NPC IDs present:
  - `9201021` Robin The Huntress
- Scripts implemented or updated:
  - `9201021.js`
- Unresolved/missing NPCs:
  - none

### `680000500` Wedding Exit Map
- Map name: `Wedding Exit map`
- Street: `Amoria`
- NPC IDs present:
  - `9201049` Ames the Wise
- Scripts implemented or updated:
  - `9201049.js`
- Unresolved/missing NPCs:
  - none

## DB Fields / Tables Changed

### `characters`
- `partnerID` nullable FK to `characters.id`
- `marriageItemID` nullable

### `marriages`
- `id`
- `husbandID`
- `wifeID`

Migration file:
- `sql/add_marriage_support.sql`

## Remaining Packet / Client TODOs

- Finish any remaining unobserved `WEDDING_ACTION` modes for v48 beyond the currently implemented registry flow.
- Finish `WEDDING_TALK` / `WEDDING_TALK_MORE` parity with full client ceremony progression.
- Reconstruct/send the correct wedding field progress packets if additional client effects are required.
- Gift registry / wishlist data is in-memory only and not persisted across restart.
- Session state is in-memory only; not crash-safe.

## Manual Test Checklist

### Engagement
- Successful engagement
- Failed engagement due to already married
- Failed engagement due to missing item
- Failed engagement due to partner offline
- Failed engagement due to partner on wrong map

### Reservation
- Reserve Cathedral wedding successfully
- Reserve Chapel wedding successfully
- Fail reservation when partner is absent
- Fail reservation when reservation ticket is missing
- Verify invitation cards are stacked

### Invitations / Guests
- Send invitation to valid guest
- Fail invitation to duplicate guest
- Open received invitation successfully
- Verify invitation cards stack correctly
- Guest entry succeeds after session starts
- Guest entry fails without valid guest ticket
- Guest enters late after ceremony has started

### Ceremony
- Wedding start with both players present
- Ceremony completion through timed stage progression
- Blessing interaction during ceremony
- Couple receives wedding rings once
- Ring/reward is not duplicated by repeated interaction

### Registry / Gifts
- Save wishlist through the client flow
- Send a gift to spouse registry
- Retrieve a gifted item once

### Persistence / Cleanup
- Relog after marriage preserves spouse state
- Disconnect during ceremony cleans up safely enough for current in-memory session design
- Divorce/reset path clears spouse state and relationship row
