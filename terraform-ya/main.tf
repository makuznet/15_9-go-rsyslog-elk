terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
}

provider "yandex" {
  token     = var.ytoken
  cloud_id  = var.ycloud
  folder_id = var.yfolder
  zone      = "ru-central1-a"
}

resource "yandex_vpc_network" "yanet" {
  name = "yanet"
}

resource "yandex_vpc_subnet" "yasubnet" {
  name           = "yasubnet"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network.yanet.id
  v4_cidr_blocks = ["192.168.8.0/24"]
}

resource "yandex_compute_instance" "ya" {
  count    = length(var.devs.prefix)
  name     = var.devs.prefix[count.index]
  hostname = var.devs.prefix[count.index]

  resources {
    cores         = 2
    memory        = 4
    core_fraction = 20
  }

  boot_disk {
    initialize_params {
      name     = "ya-${var.devs.prefix[count.index]}"
      size     = 30
      image_id = var.y_image
    }
  }

  network_interface {
    subnet_id  = yandex_vpc_subnet.yasubnet.id
    nat        = true
    ip_address = var.devs.addr[count.index]
  }

  metadata = {
    user-data = "#cloud-config\nusers:\n  - name: ${var.login_user}\n    groups: sudo\n    shell: /bin/bash\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']\n    ssh-authorized-keys:\n      - ${var.my_ssh_key}"
  }
}

resource "time_sleep" "wait_30_seconds" {

  depends_on = [yandex_compute_instance.ya, local_file.inventory]

  create_duration = "30s"
}

resource "null_resource" "known_hosts" {
  count = length(var.devs.prefix)

  provisioner "local-exec" {
    command = "ssh-keyscan -t ecdsa ${yandex_compute_instance.ya[count.index].network_interface.0.nat_ip_address} >> /Users/makuznet/.ssh/known_hosts"
  }

  depends_on = [yandex_compute_instance.ya, time_sleep.wait_30_seconds, local_file.inventory]
}

resource "local_file" "inventory" {
  content = templatefile("${path.module}/ansible_inventory.tpl",
    {
      drop_num  = range(length(var.devs.prefix))
      drop_name = var.devs.prefix
      drop_ip   = yandex_compute_instance.ya.*.network_interface.0.nat_ip_address
      drop_user = var.login_user
  })
  filename = "${path.module}/inventory.yml"

  depends_on = [yandex_compute_instance.ya]
}

output "external_ip_address" {
  value = yandex_compute_instance.ya.*.network_interface.0.nat_ip_address
}

output "internal_ip_address" {
  value = yandex_compute_instance.ya.*.network_interface.0.ip_address
}

resource "null_resource" "ansible" {

  provisioner "local-exec" {
    command = "cd ..; ansible-playbook -i terraform-ya/inventory.yml playbooks/main.yml --ask-vault-pass"
  }

  depends_on = [local_file.inventory, null_resource.known_hosts]
}
