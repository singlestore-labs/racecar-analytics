import React from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { Card } from "./components/ui/card";
import { Progress } from "./components/ui/progress";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import axios from "axios";

function App() {

  const [ecuSocketUrl] = React.useState("ws://localhost:9000/ecu/stream");
  const { lastMessage: lastEcuMessage, readyState } = useWebSocket(ecuSocketUrl);
  const [ecuData, setEcuData] = React.useState([{}]);

  const [batterySocketUrl] = React.useState("ws://localhost:9000/battery/stream");
  const { lastMessage: lastBatteryMessage } = useWebSocket(batterySocketUrl);
  const [batteryData, setBatteryData] = React.useState([{}]);

  React.useEffect(() => {
    if (lastEcuMessage !== null) {
      setEcuData((prev) => [
        ...prev.slice(-50),
        JSON.parse(lastEcuMessage.data),
      ]);
    }
  }, [lastEcuMessage]);

  React.useEffect(() => {
    if (lastBatteryMessage !== null) {
      setBatteryData((prev) => [
        ...prev.slice(-50),
        JSON.parse(lastBatteryMessage.data),
      ]);
    }
  }, [lastBatteryMessage]);

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];
  
  React.useEffect(() => {
    getAllECUData();
    getAllBatteryData();
  }, []);

  const getAllECUData = () => {
    axios.get('http://localhost:9000/ecu')
      .then(response => {
        setEcuData(response.data.slice(-40));
      })
      .catch(error => {
        console.error('Error fetching ECU data:', error);
      });
  };

  const getAllBatteryData = () => {
    axios.get('http://localhost:9000/battery')
      .then(response => {
        setBatteryData(response.data.slice(-40));
      })
      .catch(error => {
        console.error('Error fetching ECU data:', error);
      });
  };

  const Chart = ({ title, field, node }: { title: string; field: string; node: string }) => {
    return (
      <div className="mx-2 text-center">
          <h4 className="my-2">{title}</h4>
          <ResponsiveContainer width="100%" height={250}>
            <LineChart
              width={500}
              height={300}
              data={node === "ecu" ? ecuData : batteryData}
              margin={{ left: -10, right: 10 }}
            >
              <CartesianGrid strokeDasharray="4 3" stroke="#343434" />
              <XAxis
                dataKey="created_at"
                stroke="white"
                tick={{ fill: "gray" }}
                tickFormatter={(timestamp) =>
                  new Date(timestamp).toLocaleTimeString()
                }
              />
              <YAxis
                stroke="white"
                tick={{ fill: "gray" }}
                domain={["dataMin", "dataMax"]}
                allowDataOverflow={true}
              />
              <Tooltip />
              <Line
                type="monotone"
                dataKey={field}
                stroke="#820DDF"
                strokeWidth={2}
                isAnimationActive={false}
                animateNewValues={false}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
    );
  };

  return (
      <div className="p-16">
        <h1>Dashboard</h1>
        <div className="flex items-center gap-2 text-gray-400 mt-2">
          <h4>Connection Status: {connectionStatus}</h4>
          {readyState === ReadyState.OPEN && <span className="text-green-500">●</span>}
          {readyState === ReadyState.CONNECTING && <span className="text-yellow-500">●</span>}
          {(readyState === ReadyState.CLOSED || readyState === ReadyState.CLOSING || readyState === ReadyState.UNINSTANTIATED) && <span className="text-red-500">●</span>}
        </div>
        <div className="flex flex-wrap">
          <Card className="w-[150px] p-4 mr-4 mt-4">
            <div className="flex flex-col items-center justify-center m-4">
              <h4>Speed</h4>
              <h1 className="text-5xl mt-2">{ecuData[ecuData.length - 1].speed}</h1>
              <p className="text-gray-400 mt-2 font-semibold">MPH</p>
            </div>
          </Card>
          <Card className="w-[400px] p-4 mr-4 mt-4">
            <div>
              <p>Motor RPM: {ecuData[ecuData.length - 1].motor_rpm}</p>
              <Progress 
                className={`mt-2 ${
                  (ecuData[ecuData.length - 1]?.motor_rpm / 10000) * 100 < 60 ? "text-green-500" :
                  (ecuData[ecuData.length - 1]?.motor_rpm / 10000) * 100 < 80 ? "text-yellow-500" :
                  "text-red-500"
                }`}
                value={(ecuData[ecuData.length - 1]?.motor_rpm / 10000) * 100} 
              />
              <p className="mt-2">Throttle: {ecuData[ecuData.length - 1].throttle}</p>
              <Progress className="mt-2" value={ecuData[ecuData.length - 1].throttle} />
              <p className="mt-2">Brake Pressure: {ecuData[ecuData.length - 1].brake_pressure}</p>
              <Progress className="mt-2" value={(ecuData[ecuData.length - 1].brake_pressure / 10000) * 100} />
            </div>
          </Card>
          <Card className="w-[400px] p-4 mr-4 mt-4">
            <div>
              <div className="flex justify-between items-center mb-4">
                <h4 className="font-semibold">Battery</h4>
                <p className={`text-md font-semibold ${
                  batteryData[batteryData.length - 1].charge_level > 80 ? "text-green-500" :
                  batteryData[batteryData.length - 1].charge_level > 50 ? "text-yellow-500" :
                  "text-red-500"
                }`}>
                  Charge: {batteryData[batteryData.length - 1].charge_level}%
                </p>
              </div>
              <div className="flex gap-4">
                {[1, 2, 3, 4].map((cellNum) => (
                  <div key={cellNum} className="border rounded-lg p-3">
                    <h4 className="text-md font-medium mb-2">Cell {cellNum}</h4>
                    <p className="text-sm text-gray-400">
                      Voltage: {batteryData[batteryData.length - 1][`cell_voltage_${cellNum}`]} V
                    </p>
                    <p className="text-sm text-gray-400">
                      Temp: {batteryData[batteryData.length - 1][`cell_temp_${cellNum}`]} °C
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </Card>
        </div>
        <div className="flex flex-wrap">
          <Card className="w-[600px] p-4 mr-4 mt-4">
            <Chart title="Motor RPM" field="motor_rpm" node="ecu" />
          </Card>
        </div>
      </div>
  );
}

export default App;
