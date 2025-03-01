import axios from "axios";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import handleApiError from "../config/api.error";
import { apiUrl } from "../config/env.config";

type BackupDrive = {
 volume_name: string;
 drive_letter: string;
 drive_path: string;
};

const BackupDrive = () => {
 const [loading, setLoading] = useState<boolean>(false);
 const [drives, setDrives] = useState<BackupDrive[]>([]);
 const [drive, setDrive] = useState<string>("");
 const [password, setPassword] = useState<string>("");

 const handleBackup = async () => {
  if (!drive || !password) {
   //  alert("Please select a drive and enter a password");
   toast.error("Please select a drive and enter a password");
   return;
  }

  setLoading(true);

  try {
   const { data } = await axios.post(`${apiUrl}/backup-vera`, { drive_letter: drive, password });
   console.log(data);
   //  alert(data.message);
   toast.success(data.message);
  } catch (error) {
   handleApiError(error);
  } finally {
   setLoading(false);
  }
 };

 const fetchBackupDrives = async () => {
  try {
   const { data } = await axios.get(`${apiUrl}/drives`);
   setDrives(data.backup_drives);
  } catch (error) {
   handleApiError(error);
  }
 };

 useEffect(() => {
  fetchBackupDrives();
 }, []);

 return (
  <div>
   <p>Backup Drive</p>

   <select className='form-select' id='backup-drive' value={drive} onChange={(e) => setDrive(e.target.value)}>
    <option value=''>Select Drive</option>
    {drives.map((drive) => {
     return (
      <option key={drive.drive_letter} value={drive.drive_letter}>
       {drive.volume_name} ({drive.drive_letter})
      </option>
     );
    })}
   </select>

   <div className='my-1'>
    <input type='text' className='form-control' placeholder='Enter password' value={password} onChange={(e) => setPassword(e.target.value)} />
   </div>

   <div className='action-button'>
    <button className='btn btn-primary form-control' onClick={handleBackup} disabled={loading}>
     {loading ? "Backing up..." : "Backup"}
    </button>
   </div>
  </div>
 );
};

export default BackupDrive;
