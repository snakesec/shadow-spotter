rm -rf /opt/ANDRAX/kitebuilder

mkdir dist

go build -o dist/Shadow-Spotter ./cmd/shadowspotter

if [ $? -eq 0 ]
then
  # Result is OK! Just continue...
  echo "Go build Shadow-Spotter... PASS!"
else
  # houston we have a problem
  exit 1
fi

strip dist/Shadow-Spotter

if [ $? -eq 0 ]
then
  # Result is OK! Just continue...
  echo "STRIP... PASS!"
else
  # houston we have a problem
  exit 1
fi

python3 -m venv /opt/ANDRAX/kitebuilder

if [ $? -eq 0 ]
then
  # Result is OK! Just continue...
  echo "Create Virtual Environment... PASS!"
else
  # houston we have a problem
  exit 1
fi

source /opt/ANDRAX/kitebuilder/bin/activate

/opt/ANDRAX/kitebuilder/bin/pip3 install wheel

if [ $? -eq 0 ]
then
  # Result is OK! Just continue...
  echo "Pip3 install wheel... PASS!"
else
  # houston we have a problem
  exit 1
fi

/opt/ANDRAX/kitebuilder/bin/pip3 install -r UTIL-TOOLS/kitebuilder/requirements.txt

if [ $? -eq 0 ]
then
  # Result is OK! Just continue...
  echo "Pip3 install requirements... PASS!"
else
  # houston we have a problem
  exit 1
fi

cp -Rf UTIL-TOOLS/kitebuilder /opt/ANDRAX/kitebuilder/package

if [ $? -eq 0 ]
then
  # Result is OK! Just continue...
  echo "Copy PACKAGE... PASS!"
else
  # houston we have a problem
  exit 1
fi

cp -Rf dist/Shadow-Spotter /opt/ANDRAX/bin

if [ $? -eq 0 ]
then
  # Result is OK! Just continue...
  echo "Copy Shadow-Spotter bin... PASS!"
else
  # houston we have a problem
  exit 1
fi

cp -Rf andraxbin/* /opt/ANDRAX/bin
rm -rf andraxbin

chown -R andrax:andrax /opt/ANDRAX
chmod -R 755 /opt/ANDRAX