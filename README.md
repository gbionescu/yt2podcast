# yt2podcast
Turn a YouTube channel into a podcast channel using yt2podcast and listen the audio track using your favorite podcast player.

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
      }
      ```
  3. Run the application.
  4. Copy the username or channel ID of the YouTube channel that you want to listen to and add the following link in your podcast player:
  		- `http://your.domain.name/podcast/youtube/user/<username>` if you're using the username
  			- for example, CGP Grey, which can be found at `https://www.youtube.com/user/CGPGrey` would be used as `http://your.domain.name/podcast/youtube/user/CGPGrey`
  		- `http://your.domain.name/podcast/youtube/channel/<channel ID>` if you're using the channel ID
  			- the same CGP Grey example can also be used as `http://your.domain.name/podcast/youtube/channel/UC2C_jShtL725hvbm1arSV9w`

5. Wait for the application to scrape the YouTube channel. It currently uses youtube-dl to get information about the channel and it's very slow if there is a large number of uploaded videos. 
