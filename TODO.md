
This program allows one to get an idea about the state and health of the
infrastructure.

kcs-infra   --output table,json,csv,text --environment default=dit
  db
  backup
    list
    check
  instance
    list --include-stopped, --filter name,--ssm, --no-ssm, 
    connect NAME or ID
    describe NAME or ID
  app-config
    list --filter
    update NAME --value
    revert NAME
    
    
    
Future
  logs
  metrics
  
  
Other Apps


for releases:
chez moi - has an execlent golintvetbuild for binary release of go 
