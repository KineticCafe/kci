
This program allows one to get an idea about the state and health of the
infrastructure.


kci

- [X] instance - commands for dealing with instances
  - [X] list - list instances, with filters for status and name
  - [X] aging - list old instances and images
  - [X] ssm - list instances that are SSM managed
  - [X] reboot - reboot a single instance
  - [X] scan - report on OS, reboot status, security update count
- [ ] ssm - commands for dealing with SSM
  - [X] list - an alias of `instance ssm`
  - [ ] session - open a session (ssh without bastion)
  - [ ] patch - list available patches
  - [ ] update - apply patches
  - [ ] run - run either a script or a command
- [ ] param - commands for SSM parameter store 
- [ ] database - commands for databases
  - [ ] list - list all databases
- [ ] snapshot - commands for database snapshots
 - [ ] create - create a snapshot of an RDS instance
 - [ ] list - list all snapshots
 - [ ] test - test a snapshot by creating a new database
 - [ ] restore - restore snapshot to a given RDS instance
- [ ] route - Route53 stuff
- [ ] iam - user and identity management
  - [ ] disable - disable a user
  - [ ] list - list all users, filters for console and key
  - [ ] aging - an aging list
  - [ ] rotate - rotate a single key


## For the Future

- Patch Compliance: use SSM to check and manage patch compliance for your EC2 instances. 
- Automated Remediation: SSM allows you to automatically remediate non-compliant resources through 
  - AWS Config and SSM Automation documents
- Vulnerability Scanning: SSM integrates with AWS Inspector to provide automated security assessment service. 
- Intrusion Detection: automate the deployment of intrusion detection and prevention systems. 


### database

- Create Snapshot: create a manual snapshot of a specified RDS database instance.
  - kci rds snapshot create --db-instance mydb
- List Snapshots: list all available manual and automatic snapshots for a specified RDS database instance.
  - kci rds snapshot list --db-instance mydb
- Restore Database: This command would restore a database from a specified snapshot, creating a new RDS database instance.
  - kci rds snapshot restore --snapshot-id mysnapshot --db-instance
- Test Restore Procedure: restore a database from a snapshot to new instance, run test script, clobber new id


