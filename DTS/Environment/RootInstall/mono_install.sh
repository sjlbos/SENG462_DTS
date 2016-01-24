# This script installs the Mono framework, used to run .NET applications in a unix/linux environment

echo "Adding Mono Project signing key..."
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF

echo "Adding Mono package repository..."
echo "deb http://download.mono-project.com/repo/debian wheezy main" | tee /etc/apt/sources.list.d/mono-xamarin.list

echo "Updating package list..."
apt-get update

echo "Installing package \"mono-complete\""
apt-get -y install mono-complete 