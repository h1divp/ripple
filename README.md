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

Dev slice-of-life
- versioned deployments
- good observability
- sharded db
- testing :)

UI
- i would like it to be a bit more stylized than just boilerplate, but that's not very important early on
