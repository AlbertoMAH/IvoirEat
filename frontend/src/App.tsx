import { Refine } from "@refinedev/core";
import dataProvider from "@refinedev/simple-rest";

function App() {
  return (
    <Refine
      dataProvider={dataProvider("https://api.fake-rest.refine.dev")}
    >
      <div>Ceci est un test minimal de Refine. Si ce message s'affiche, le problème vient du routeur ou d'AntD.</div>
    </Refine>
  );
}

export default App;
