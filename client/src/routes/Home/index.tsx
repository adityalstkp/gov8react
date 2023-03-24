import {
  containerStyle,
  linkStyle,
  subTitleStyle,
  wordStyle,
} from "../../styles/common";

const Home = () => {
  const greet = "";
  return (
    <div className={containerStyle}>
      <h1 className={wordStyle}>gov8react</h1>
      {greet && (
        <>
          <h3 className={subTitleStyle}>Your User Agent:</h3>
          <p className={wordStyle}>{greet}</p>
        </>
      )}
      <a className={linkStyle} href="/about">
        About
      </a>
    </div>
  );
};

export default Home;
