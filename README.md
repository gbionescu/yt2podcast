# yt2podcast
Turn a YouTube channel into a podcast channel using yt2podcast and listen audio tracks using your favorite podcast player.

Tested with:
  - Podcast Addict on Android (works)
  - gPodder on Linux (works)
  - iTunes on Mac (does not work) - probably the XML needs more tags to satisfy iTunes

How to use:
  1. Build the package.
  2. Create a json file named `config.json` that specifies how podcast clients connect to the podcast server:
      ```
      {
        "Hostname": "your.domain.name",
        "port": "8080"
        "MaxYTStorageBytes": "How much data to keep in the YouTube storage folder. Value set in bytes."
      }
      ```
  3. Get YouTube API credentials from your GCP account and save them to `client_secret.json`.
  4. Run the application.
  5. Copy the username or channel ID of the YouTube channel that you want to listen to and add the following link in your podcast player:
  		- `http://your.domain.name:port/api/ytchan/<username or channel ID>`
  		- for example, CGP Grey, which can be found at `https://www.youtube.com/user/CGPGrey` would be used as `http://your.domain.name:port/api/ytchan/CGPGrey` or `http://your.domain.name:port/api/ytchan/UC2C_jShtL725hvbm1arSV9w`
  6. Note: You will need to authenticate your application with the YouTube servers after adding a channel or playlist, so keep an eye on the console where the application was launched.
