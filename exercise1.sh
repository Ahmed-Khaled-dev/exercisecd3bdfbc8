#!/bin/bash

readonly FILESYSTEM_IMAGE="helloworld.img"
readonly KERNEL_NAME="ubuntu-vmlinuz"
readonly PREBUILT_KERNEL_LINK="https://cloud-images.ubuntu.com/minimal/releases/jammy/release/unpacked/ubuntu-22.04-minimal-cloudimg-amd64-vmlinuz-generic"

# Install dependencies
command -v qemu-system-x86_64 >/dev/null || sudo apt install -y qemu-system-x86
command -v wget >/dev/null || sudo apt install -y wget

# Create filesystem image
dd if=/dev/zero of=$FILESYSTEM_IMAGE bs=1M count=8
mkfs.ext4 -F $FILESYSTEM_IMAGE

# Mount and create rootfs
mkdir -p mnt
sudo mount $FILESYSTEM_IMAGE mnt

# Essential directory structure
sudo mkdir -p mnt/{bin,dev}

# Install busybox and use busybox's sh
sudo apt install -y busybox-static
sudo cp /bin/busybox mnt/bin/
sudo ln -s busybox mnt/bin/sh

# Init script
echo '#!/bin/sh
echo "hello world"
exec /bin/sh' | sudo tee mnt/init >/dev/null
sudo chmod +x mnt/init

# Finished creating rootfs
sudo umount mnt
rmdir mnt

# Download prebuilt kernel
if [ ! -f "$KERNEL_NAME" ]; then
    wget -O $KERNEL_NAME $PREBUILT_KERNEL_LINK
fi

# Run QEMU with proper parameters
qemu-system-x86_64 \
    -kernel "$KERNEL_NAME" \
    -append "console=ttyS0 root=/dev/sda rw init=/init noreplace-smp" \
    -drive "file=$FILESYSTEM_IMAGE,format=raw,if=ide" \
    -nographic \
    -m 512M