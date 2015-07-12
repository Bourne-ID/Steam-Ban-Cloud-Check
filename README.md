# Steam-Ban-Cloud-Check
Checks for steam bans using Steam API and Google App Cloud

Created to enhance the user experience by removing the step to use their own API key (which can be problematic if they do not own a domain name). As per the Steam API Terms of Use, API keys must be kept confidential, hence cannot be embedded into a Chrome Extension, and therefore the creation of this cloud service.

##Input
```
post data: steamids=xxx,yyy,zzz
```
##Output
JSON Example:
```
[{"SteamID":"76561197962361621","CommunityBanned":false,"VACBanned":false,"NumberOfVACBans":0,"DaysSinceLastBan":0,"NumberOfGameBans":0,"EconomyBan":"none"},{"SteamID":"76561197960434622","CommunityBanned":false,"VACBanned":false,"NumberOfVACBans":0,"DaysSinceLastBan":0,"NumberOfGameBans":0,"EconomyBan":"none"}]
```

##Personal Note
This is my first venture into Google App Cloud coding, and GO language in general so I suspect there are better ways of doing many things here, along with ways to make the code better. As such, you probably shouldn't use this as an example of best practices for GO Code.

##Future editions plans:

v2: Use Datastore to cache common user searches (reduces stress of Steam API Key)

v3: Anonymous Usage statistics on API calls (chrome extension opt-in tickbox)

v4: (opt-in on chrome extension) Chrome User Watch List, automated monitoring and notification pushing to client
