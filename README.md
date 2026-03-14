# Echo v2
Geo-based real-time messaging. Chat with everyone in lecture or at your local library.

This repository is currently not open for contribution, but contributors may be accepted in the future if this project gets off the ground (!)

## How it works

## Hosting

## History

## MVP
- chat room with main functionality
- settings page/menu pop-up
- optional auth
  - allows for completely anonymous messaging
  - allows for quick onboarding
  - ip should still be stored in db for moderation purposes (temp ban)
  - users who want to want to chat showing reputation, keep dms etc, etc can sign up

Nice to haves
- reactions
- Dice roll-ed profile picture & user name
- Map
  - shows exactly where message reach is
- Signposts (come up with better name)
  - Users can leave a note on a signpost (typed or maybe even drawn)
- Distributed mode
  - Running instance is in charge of a geographically scoped area and connects to larger network
- Moderator tools
  - use anonymized reports
  - shadowbanning or reputation score
    - reputation score would incentivize users to keep accounts for longer
      - better reputation -> can interact more? (i.e. have larger reach, more possible reactions, social decorations)
  - network wide announcements
  - DMs
    - but to make it more interesting, you can only DM a user and see their profile while they're online
    - geographically unrestricted of course (bad privacy issue otherwise)
- bluetooth only mode
  - when resources within an area are highly utilized a this mode can be suggested as a fallback
  - likely only allows for a very small radius and may be hard to do in a web app; needs investigation
- phone number sign-up
- user reputation
  - calculated based on things like amount of messages sent (activity), friends, reacted messages
  - could use up/down voting?
  - "shadowmuting" idea:
    - user has been downvoted so frequently in a short period of time that their messages start being muted without the sendering knowing it was fully sent or not
    - natural way to deter hatespeech and bots without needing to flag things for moderation
    - can make it obvious too, i.e. disable message box for time period, pop up ban message, show less messages to user
- add something reminiscent of 3ds streetpass 

Dev slice-of-life
- versioned deployments
- good observability
- sharded db
- testing :)

UI
- i would like it to be a bit more stylized than just boilerplate, but that's not very important early on
- opengraph
- visual indictator (like 0.1s blink) for when the user's location is updated
