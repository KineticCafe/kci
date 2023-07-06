
This program allows one to get an idea about the state and health of the
infrastructure.


kci

- [X] instance - commands for dealing with instances
  - [X] list - list instances, with filters for status and name
  - [X] aging - list old instances and images
  - [X] ssm - list instances that are SSM managed
  - [X] reboot - reboot a single instance
  - [X] scan - report on OS, reboot status, security update count
- [.] ssm - commands for dealing with SSM
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
 - [ ] test - test a snapshot by createing a new database
- [ ] route - Route53 stuff
- [ ] iam - user and identity management
  - [ ] disable - disable a user
  - [ ] list - list all users, filters for console and key
  - [ ] aging - an aging list
  - [ ] rotate - rotate a single key


## For the Future

Patch Compliance: As mentioned earlier, you can use SSM to check and manage patch
compliance for your EC2 instances. You can create a command that checks for missing
patches and another command to apply them.

Configuration Compliance: SSM allows you to define and enforce system configurations
using State Manager. You could include a command in your CLI tool that reports on
instances that are out of compliance.

Automated Remediation: SSM also allows you to automatically remediate non-compliant
resources through AWS Config and SSM Automation documents. You could potentially
include a command to trigger these remediations manually.

Secret Rotation: If you're using SSM Parameter Store to manage secrets, you could
include a command to rotate secrets.

Vulnerability Scanning: SSM integrates with AWS Inspector to provide automated
security assessment service. This helps to improve the security and compliance of
applications deployed on AWS. CLI command could be created to run assessments and
report on findings.

Intrusion Detection: With SSM, you can automate the deployment of intrusion detection
and prevention systems. While this might not lend itself to a CLI command, it could
be part of your security setup.


### database

Create Snapshot: This command would create a manual snapshot of a specified RDS
database instance.

Example usage: myapp rds snapshot create --db-instance mydb

List Snapshots: This command would list all available manual and automatic snapshots
for a specified RDS database instance.

Example usage: myapp rds snapshot list --db-instance mydb

Delete Snapshot: This command would delete a specified snapshot.

Example usage: myapp rds snapshot delete --snapshot-id mysnapshot

Restore Database: This command would restore a database from a specified snapshot,
creating a new RDS database instance.

Example usage: myapp rds snapshot restore --snapshot-id mysnapshot --db-instance
newdb

Test Restore Procedure: This command would restore a database from a snapshot, run a
set of user-specified tests to ensure the restored database is functioning correctly
(this might involve running a script or set of SQL commands), and then delete the
restored database instance.

Example usage: myapp rds snapshot test-restore --snapshot-id mysnapshot --test-script
mytest.sql



