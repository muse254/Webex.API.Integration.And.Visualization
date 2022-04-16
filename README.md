# Integration and Visualization with Webex API

## Description

Through integration with Webex's API visualize meeting quality afterwards. Successful execution and description will show that the student has a solid understanding of the visualization of data through the Rest API.
Knowledge of REST API and Oauth is an advantage. The assignment consists of the following:
- Register an Integration with the API https://developer.webex.com/docs/integrations
- Create scripts that retrieve meeting room quality (in JSON format) and save it
- Visualize the meeting room experience, a starting point may be what you see in Analytics in Webex Control Hub, but you are free to be creative and create alternatives

## Task

1. Create a Web server as the test environment for the integration. 
2. Register an Integration with the API https://developer.webex.com/docs/integrations
3. Determine which Scope is needed (minimum: analytics:read_all) 
4. When creating integration, you will get the following that must be taken care of: a. Client ID 
    - Client Secret 
    - OAuth Authorization URL (NOTE: each time you add  new scopes this will  update) 
    - Integration ID 
5. Create scripts that safeguard Token - and optionally update it if you wish. 
6. Create scripts that collect the right ID from the meeting(s) you want to visualize.  The correct id must be in the format {string_|_string}. 
 In advance, you often only  have access code, ex: 23726108570 
Create scripts that retrieve meeting room quality (in json format) and save it (one can only make one call per 5 min.)   
7. Visualizing the meeting room experience, a starting point may be what you see in Analytics in the Webex Control Hub, but you are free to be creative and create alternatives.  

A good description of how to integrate with the Webex API can be found here: https://developer.webex.com/blog/real-world-walkthrough-of-building-an-oauth-webexintegration 
 
There are two pointers to Webex's api: 
- https://webexapis.com (token management, users, meeting overviews, etc.). 
- https://analytics.webexapis.com (meeting room activity, quality and overview). 


## REST API document as jon: 

https://developer.webex.com/docs/api/getting-started (with api references In left menu) 
 
## Tips: 

Take API calls directly from server, and not from browser (to avoid CORS issues).	


