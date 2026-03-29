# Echo v2
Geo-based real-time messaging. Chat with everyone in lecture or at your local library or cafe.

This repository is currently not open for contribution, but contributors may be accepted in the future if this project gets off the ground (!)

## How it works

## Hosting

## History

## MVP
- chat room with main functionality
  - TODO:
    - rate limiting: message length on both api/web, fast message sending
    - disconnect user if location cannot be retrieved, show error message
- settings page/menu pop-up
- add TTL to location entries
  - i dont think this is a good idea for the sessions since we dont want to desync from the state of the websocket connections
  - if geoindex isnt found an update could be requested from the client, if that fails then we need to show an error in UI
- optional auth
  - allows for completely anonymous messaging
  - allows for quick onboarding
  - ip should still be stored in db for moderation purposes (temp ban)
  - users who want to want to chat showing reputation, keep dms etc, etc can sign up
- trust
  - each session should have a trust score?
  - if it goes negative past a certain threshhold, various things can happen
    - a user cannot type for an allotted period of time
    - user will have their session kicked
  - !! look up best practices for this, should not be able to be gamed by botshttps://marcoislandchamber.org/marco-event/pancake-breakfast-3/
- for transparency, add pop up modal when app loads asking for location access
  - do so before the user obtains a session, because we can use the session an indicator that they agreed to use location (along with browser permission)

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
    - actually probably a bad idea, since botted users could spam downvoted on someone. ill leave this idea here for the meantime though
  - "shadowmuting" idea:
    - user has been downvoted so frequently in a short period of time that their messages start being muted without the sendering knowing it was fully sent or not
    - natural way to deter hatespeech and bots without needing to flag things for moderation
    - can make it obvious too, i.e. disable message box for time period, pop up ban message, show less messages to user
    - definitely dont want to use this in the wrong way, again look up best practices
- add something reminiscent of 3ds streetpass 
- when there are no users present in someones area, try the following
  - allow a user to expand their radius to a larger present (e.g. 50m, 200m, 1km, etc).
  - only let users talk to eachother within the same radius, but also who have the *same radius setting*
  - this allows for a higher chance of interaction (making the app more enjoyable) while still keeping some locality. i.e. people in the same city can talk to eachother if theres no-one in the same building in chat
- try out terraform
- jp/eng localization, check user agent for default language
- emoji picker with :aliases:

Dev slice-of-life
- versioned deployments
- good observability
- sharded db
- testing :)
  - would be nice to have a script for load testing ws connections and flood of chat messages

UI
- i would like it to be a bit more stylized than just boilerplate, but that's not very important early on
- opengraph
- visual indictator (like 0.1s blink) for when the user's location is updated
- (on desktop) drag-to-resize chatbox
- clicking on the logo at the top should show a popup with an explanation, maybe link to blogpost.
