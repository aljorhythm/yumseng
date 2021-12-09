import { useEffect } from "react";
import ReactPlayer from "react-player";
import src from "../assets/yumseng-full.mp3";
interface CheerVolumeBarProp {
  intensity: number;
}
const generateBar = (strength: number): string => {
  let result = "";
  while (strength--) {
    result += "-";
  }
  return result;
};
const totalStrength = 7;
const standardizeVolume = (strength: number): number =>
  strength / totalStrength;
const url = `yumseng-full.mp3`;
const CheerVolumeBar = ({ intensity }: CheerVolumeBarProp) => {
  useEffect(() => {
    console.log("is url playable?", ReactPlayer.canPlay(url));
  }, []);
  return (
    <>
      <ReactPlayer
        url={src}
        playing
        volume={standardizeVolume(intensity)}
        width="100%"
        height="100%"
      />
      <div>{standardizeVolume(intensity)}</div>
      <div>
        <div>{generateBar(intensity)}</div>
      </div>
    </>
  );
};

export default CheerVolumeBar;
