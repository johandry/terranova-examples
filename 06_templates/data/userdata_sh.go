package data

// UserdataTmpl is the bash code in form of Go template to set as USERDATA and
// execute by Cloud Init in the new instance after rendered.
const UserdataTmpl = `#!/bin/bash

# Start the status reporter
html_status() { 
	echo '<li><span style="color:blue">[INFO]</span> '$1'</li>' >> status.txt; 
	echo "<html><head><title>Server Status</title>" > index.html;
	[[ $2 != "END" ]] && echo '<meta http-equiv="refresh" content="10"/>' >> index.html; 
	echo '</head><body><h2>Status:</h2><ul style="list-style-type:none">' >> index.html; 
	cat status.txt >>index.html;
	echo '</ul><h2>Logs:</h2><textarea rows="50" cols="100">' >>index.html
	cat /var/log/cloud-init-output.log >>index.html
	echo '</textarea>' >>index.html
	echo "</body></html>" >>index.html; 
}
%{ if status }
nohup busybox httpd -f -p ${status_port} &
echo $! > httpd.pid
%{ endif }

# Install Docker
html_status "updating the packages..."
sudo apt-get update
html_status "done"
html_status "installing curl, CA certificates and other packages..."
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
html_status "done"
html_status "installing docker..."
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo usermod -aG docker $USER
html_status "done"

# Install Docker Compose
html_status "installing docker compose..."
sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
html_status "done"

# Create the Docker Compose file
html_status "creating docker compose file..."
echo '${ docker_compose_b64 }' | base64 --decode > docker-compose.yaml
html_status "done"

html_status "starting docker compose with Let'sChat and MongoDB..."

# Start Docker Compose
sudo docker-compose up -d
html_status "done"

TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/public-hostname)

html_status 'Go to: <a href="http://'$IP':${ letschat_port }">http://'$IP':${ letschat_port }</a>' "END"
%{ if status }
# TODO: Make this a task that will run 10 min later
# kill -9 $(cat httpd.pid)
# rm -f httpd.pid
%{ endif }
`
