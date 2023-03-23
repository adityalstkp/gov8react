import React from "react";
import { containerStyle, subTitleStyle, wordStyle } from "../../styles/common";

interface HomeProps {
  data: Record<string, unknown>;
}

const Home = (props: HomeProps) => {
  const greet = props.data.greet as string;
  return (
    <div className={containerStyle}>
      <h1 className={wordStyle}>gov8react</h1>
      {greet && (
        <>
          <h3 className={subTitleStyle}>Your User Agent:</h3>
          <p className={wordStyle}>{greet}</p>
        </>
      )}
    </div>
  );
};

export default Home;
