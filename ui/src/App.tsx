import axios from "axios";
import { useState } from "react";
import { toast, ToastContainer } from "react-toastify";
import Action from "./components/Action";
import BackupDrive from "./components/backup_drive";
import MountDrive from "./components/GenerateDrive";
import handleApiError from "./config/api.error";
import { apiUrl } from "./config/env.config";

const GeneratePassword = () => {
 const [password, setPassword] = useState<string>("");
 const [encryptedPassword, setEncryptedPassword] = useState<string>("");

 const handleGeneratePassword = async () => {
  try {
   const { data } = await axios.post(`${apiUrl}/encrypt-password`, { password });
   setEncryptedPassword(data.password);
   if (navigator.clipboard) {
    try {
     await navigator.clipboard.writeText(data.password);
     toast.success("Password copied to clipboard");
    } catch (error) {
     console.log(error);
     toast.error("Failed to copy password to clipboard");
    }
   }
   console.log(data);
  } catch (error) {
   handleApiError(error);
  }
 };

 return (
  <div className='mb-3 mt-4'>
   <p>Generate Password</p>

   <div className='input-group mb-3'>
    <input type='password' className='form-control' placeholder='password' value={password} onChange={(e) => setPassword(e.target.value)} />
   </div>

   <p>
    Encrypted Password: <strong>{encryptedPassword}</strong>
   </p>

   <button className='btn btn-primary' onClick={handleGeneratePassword}>
    generate
   </button>
  </div>
 );
};

// const DriveList = () => {
//  const [drives, setDrives] = useState<string[]>([]);
//  const fetchDrives = async () => {
//   try {
//    const { data } = await axios.get(`${apiUrl}/mounted-drives`);
//    console.log(data);
//   } catch (error) {
//    handleApiError(error);
//   }
//  };

//  useEffect(() => {
//   fetchDrives();
//  }, []);

//  return (
//   <div onClick={() => {}}>
//    <p>Drive List</p>
//   </div>
//  );
// };

function App() {
 return (
  <>
   <ToastContainer />
   <div className='container mt-4' style={{ width: "400px" }}>
    {/* <DriveList /> */}
    <MountDrive />

    <Action />
    <GeneratePassword />

    <BackupDrive />
   </div>
  </>
 );
}

export default App;
