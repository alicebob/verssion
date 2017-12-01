# helper to bootstrap some content
for l in $(cat pages.txt); do
    echo $l
    curl -s -S http://localhost:3141/p/$l/ > /dev/null
    sleep 5
done
