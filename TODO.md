
This program allows one to get an idea about the state and health of the
infrastructure.

kci   --output table,json,csv,text --environment default=dit
  db
  backup
    list
    check
  route
    describe URL
  
  
Jobs I want to do...

kci
  instance
    public - open ports, public zone
    status - OS, needs reboot, etc.
  dns
    list --zone --filter - lists all records across all (or --zone) zones
    zone ZONEID - describes a zone or list if not there
    url URL - describes the record, including zone and what it attaches to
  backup
    dump - checks on status of database dumps
    snapshot - checks on status of database snapshots
  ssm-connect
    enable --user
    disable --user
  iam
    rotate-access-key --user - returns the access key and secret for sending
    disable --user --no-console --no-access-key
    list --ageing --filter




    
    
    
Future
  logs
  metrics
  
  
Other Apps


for releases:
chez moi - has an execlent golintvetbuild for binary release of go 
