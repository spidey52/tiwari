import axios from "axios";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import handleApiError from "../config/api.error";
import { apiUrl } from "../config/env.config";

const fetchDrives = async () => {
 try {
  const { data } = await axios.get(`${apiUrl}/drives`);
  return data.drives;
 } catch (error) {
  console.error(error);
  return [];
 }
};

export type Drive = {
 volume_name: string;
 drive_letter: string;
 drive_path: string;
};

const MountDrive = () => {
 const [loading, setLoading] = useState<boolean>(false);
 const [unmounting, setUnmounting] = useState<boolean>(false);
 const [password, setPassword] = useState<string>("");
 const [drive, setDrive] = useState<string>("");
 const [drives, setDrives] = useState<Drive[]>([]);

 useEffect(() => {
  fetchDrives().then((drives) => setDrives(drives));
 }, []);

 const handleMount = async () => {
  if (!drive || !password) {
   toast.error("Please select a drive and enter a password");
   return;
  }

  setLoading(true);

  try {
   const { data } = await axios.post(`${apiUrl}/mount`, { drive_letter: drive, password });
   toast.success(data.message);
  } catch (error) {
   handleApiError(error);
  } finally {
   setLoading(false);
  }
 };

 const handleUnmount = async () => {
  if (!drive) {
   toast.error("Please select a drive");
   return;
  }

  setUnmounting(true);

  try {
   const { data } = await axios.post(`${apiUrl}/unmount`, { drive_letter: drive });
   toast.success(data.message);
  } catch (error) {
   handleApiError(error);
  } finally {
   setUnmounting(false);
  }
 };

 return (
  <>
   <select className='form-select mb-3' value={drive} onChange={(e) => setDrive(e.target.value)}>
    <option selected>Select Drive</option>

    {drives.map((drive) => (
     <option key={drive.drive_letter} value={drive.drive_letter}>
      {drive.volume_name}
     </option>
    ))}
   </select>
   <div className='mb-3'>
    <input type='password' className='form-control' value={password} onChange={(e) => setPassword(e.target.value)} placeholder='dummy-password' />
   </div>

   <div className='mb-3'>
    <button className='btn btn-success form-control' onClick={handleMount} disabled={loading}>
     mount drive
    </button>
    <p></p>
    <button className='btn btn-danger form-control' onClick={handleUnmount} disabled={unmounting}>
     unmount drive
    </button>
   </div>
  </>
 );
};

export default MountDrive;
