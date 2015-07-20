all:
	go build .
	mkdir -p build
	echo 2.0 > build/debian-binary
	echo "Package: sentrynotifier" > build/control
	echo "Version: 0.0.1" >> build/control
	echo "Architecture: amd64" >> build/control
	echo "Section: net" >> build/control
	echo "Installed-Size: `du -k sentrynotifier | cut -f1`" >> build/control
	echo "Maintainer: Will Barrett <will@speak.io>" >> build/control
	echo "Priority: optional" >> build/control
	echo "Description: A command line tool for reporting errors to Sentry" >> build/control
	echo " Built" `date`
	sudo rm -rf build/usr
	mkdir -p build/usr/local/bin
	cp sentrynotifier build/usr/local/bin/sentrynotifier
	sudo chown -R root: build/usr
	tar cvzf build/data.tar.gz -C build usr
	tar cvzf build/control.tar.gz -C build control
	cd build && ar rc sentrynotifier.deb debian-binary control.tar.gz data.tar.gz && cd ..
