import { useIntro } from "../../hooks/intro";
import { containerStyle, linkStyle, wordStyle } from "../../styles/common";

const Home = () => {
  const { data: greet } = useIntro();

  return (
    <div className={containerStyle}>
      <h1 className={wordStyle}>gov8react</h1>
      {greet && <p className={wordStyle}>{greet}</p>}
      <a className={linkStyle} href="/about">
        About
      </a>
    </div>
  );
};

export default Home;
