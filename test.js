(async () => {
  let promises = [];

  for (i = 0; i < 7000; i++) {
    const data = fetch("http://localhost:3000/auction", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "forsen",
      },
      body: JSON.stringify({
        robloxId: Number((Math.random() * 1000000).toFixed(0)),
        robloxName: "XQCSLAMTHEFART",
        itemType: "PET",
        startPrice: 100000,
        itemData: {
          e: false,
          id: "2",
          nk: "Cat",
          xp: 0,
          idt: 1763,
          lvl: 1,
          uid: "ide42e4a8a5cb745e6bc89a67778d8b144",
          place: 0,
        },
      }),
    }).then((res) => res.json());

    promises.push(data);
    // promises.push(data);
    // promises.push(data);
    // const playtime = await fetch("https://v3.kattah.me/leaderboard/playtime", {
    //   method: "GET",
    //   headers: {
    //     "Content-Type": "application/json",
    //     Authorization: "2b9011ccf6d803a3",
    //   },
    // }).then((res) => res.json());
    // console.log(playtime.data.nof2p, "playtime");
    // const robux = await fetch("https://v3.kattah.me/leaderboard/robux", {
    //   method: "GET",
    //   headers: {
    //     "Content-Type": "application/json",
    //     Authorization: "2b9011ccf6d803a3",
    //   },
    // }).then((res) => res.json());
    // console.log(robux, "robux");
    // const power = await fetch("https://v3.kattah.me/leaderboard/power", {
    //   method: "GET",
    //   headers: {
    //     "Content-Type": "application/json",
    //     Authorization: "2b9011ccf6d803a3",
    //   },
    // }).then((res) => res.json());
    // console.log(power, "power");
    // const eggs = await fetch("https://v3.kattah.me/leaderboard/eggs", {
    //   method: "GET",
    //   headers: {
    //     "Content-Type": "application/json",
    //     Authorization: "2b9011ccf6d803a3",
    //   },
    // }).then((res) => res.json());
    // console.log(eggs, "eggs");
    // const secrets = await fetch("https://v3.kattah.me/leaderboard/secrets", {
    //   method: "GET",
    //   headers: {
    //     "Content-Type": "application/json",
    //     Authorization: "2b9011ccf6d803a3",
    //   },
    // }).then((res) => res.json());
    // console.log(secrets, "secrets");
  }

  const data = await Promise.all(promises);
  console.log(data);
})();
