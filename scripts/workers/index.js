addEventListener("fetch", (event) => {
  event.respondWith(handleRequest(event.request));
});

async function handleRequest(request) {
  let pathname = new URL(request.url).pathname;
  if (pathname === "/subscribe") {
    //TODO check HTTP Request type
    try {
      const body = await request.text();
      let date = Date.now();
      //TODO perform email validation
      //TODO check if email already exists
      //TODO send a confirmation email before actually adding to DB
      await NL_EMAIL.put(body, date);
      return new Response("Subscribed", { status: 201 });
    } catch (error) {
      return new Response(error || "error", { status: 400 });
    }
  } else if (pathname === "/unsubscribe") {
    try {
      const body = await request.text();
      //TODO maybe use encoded request and decode it here
      await NL_EMAIL.delete(body);
      return new Response("Unubscribed", { status: 204 });
    } catch (error) {
      return new Response(error || "error", { status: 400 });
    }
  }
  const value = await NL_EMAIL.list();
  return new Response(JSON.stringify(value.keys));
}
