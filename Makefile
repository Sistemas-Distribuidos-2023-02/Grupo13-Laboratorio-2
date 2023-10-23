docker-latinoamerica:
	docker-compose -f docker-compose-latinoamerica.yml build
	docker-compose -f docker-compose-latinoamerica.yml up --remove-orphans

docker-europa:
	docker-compose -f docker-compose-europa.yml build
	docker-compose -f docker-compose-europa.yml up --remove-orphans

docker-asia:
	docker-compose -f docker-compose-asia.yml build
	docker-compose -f docker-compose-asia.yml up --remove-orphans

docker-australia:
	docker-compose -f docker-compose-australia.yml build
	docker-compose -f docker-compose-australia.yml up --remove-orphans

docker-oms:
	docker-compose -f docker-compose-oms.yml build
	docker-compose -f docker-compose-oms.yml up --remove-orphans

docker-onu:
	docker-compose -f docker-compose-onu.yml build
	docker-compose -f docker-compose-onu.yml up --remove-orphans

docker-dn1:
	docker-compose -f docker-compose-dn1.yml build
	docker-compose -f docker-compose-dn1.yml up --remove-orphans

docker-dn2:
	docker-compose -f docker-compose-dn2.yml build
	docker-compose -f docker-compose-dn2.yml up --remove-orphans

docker-vm049:
	docker-compose -f docker-compose-vm049.yml build
	docker-compose -f docker-compose-vm049.yml up --remove-orphans

docker-vm050:
	docker-compose -f docker-compose-vm050.yml build
	docker-compose -f docker-compose-vm050.yml up --remove-orphans

docker-vm051:
	docker-compose -f docker-compose-vm051.yml build
	docker-compose -f docker-compose-vm051.yml up --remove-orphans

docker-vm052:
	docker-compose -f docker-compose-vm052.yml build
	docker-compose -f docker-compose-vm052.yml up --remove-orphans
