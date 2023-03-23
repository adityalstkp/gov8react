// for testing, benchmarking and debugging only
import express from "express";
import { getMatchRoutes, renderApp } from "./app";

const app = express();
app.get("*", async (req, res) => {
  const match = getMatchRoutes(req.url);
  if (!match) {
    res.writeHead(404);
    return;
  }

  console.log(document);
  const data = await renderApp({ url: req.url, staticData: null });
  res.setHeader("content-type", "text/html");
  res.send(`
<html>
<div id="root">${data}</data>
</html>
`);
});

app.listen(3000, () => {
  console.log("node express listening");
});
