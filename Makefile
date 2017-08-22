.PHONY: all install

INSTALL_DIR := /usr/lib/go-luks-suspend

all: go-luks-suspend initramfs-suspend

go-luks-suspend:
	GOPATH="$(CURDIR)" go build goLuksSuspend/cmd/go-luks-suspend

initramfs-suspend:
	GOPATH="$(CURDIR)" go build goLuksSuspend/cmd/initramfs-suspend

install: all
	install -Dm755 go-luks-suspend "$(DESTDIR)$(INSTALL_DIR)/go-luks-suspend"
	install -Dm755 initramfs-suspend "$(DESTDIR)$(INSTALL_DIR)/initramfs-suspend"
	install -Dm644 initcpio-hook "$(DESTDIR)/usr/lib/initcpio/install/suspend"
	install -Dm644 systemd-suspend.service "$(DESTDIR)/etc/systemd/system/systemd-suspend.service"

clean:
	rm -f go-luks-suspend initramfs-suspend

# vim:set sw=4 ts=4 noet:
