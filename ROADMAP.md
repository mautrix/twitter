# Features & roadmap

* Matrix → Twitter
  * [ ] Message content
    * [x] Text
    * [ ] Formatting
    * [ ] Media
      * [x] Images
      * [x] Videos
      * [x] Gifs
  * [x] Message reactions
  * [ ] Typing notifications
  * [x] Read receipts
* Twitter → Matrix
  * [ ] Message content
    * [x] Text
    * [ ] Formatting
    * [x] Media
      * [x] Images
      * [x] Videos
      * [x] Gifs
  * [x] Message reactions
  * [x] Message history
    * [x] When creating portal
    * [x] Missed messages
  * [x] Avatars
  * [ ] † Typing notifications
  * [ ] † Read receipts
* Misc
  * [x] Automatic portal creation
    * [x] At startup
    * [x] When receiving invite or message
  * [x] Provisioning API for logging in
  * [ ] Private chat creation by inviting Matrix puppet of Twitter user to new room
  * [ ] Option to use own Matrix account for messages sent from other Twitter clients
    * [x] Automatic login with shared secret
    * [ ] Manual login with `login-matrix`
  * [x] E2EE in Matrix rooms

† Information not automatically sent from source, i.e. implementation may not be possible
‡ Maybe, i.e. this feature may or may not be implemented at some point
