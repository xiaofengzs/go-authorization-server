1. build login function
   -  if has session in redis, redirect user to index
   -  if no session, redirect user to login
   -  if user try to login, no sessoin in cookie and redis, check username

2. improve login
   - if login is from redirect, should go back to origin url - finished