import axios from "axios";
import handleApiError from "../config/api.error";
import { apiUrl } from "../config/env.config";
import { toast } from "react-toastify";

const Action = () => {
 const handleShutdown = async () => {
  try {
   const { data } = await axios.post(`${apiUrl}/shutdown`);
  //  alert(data.message);
  toast.success(data.message);
  } catch (error) {
   handleApiError(error);
  }
 };
 const handleReboot = async () => {
  try {
   const { data } = await axios.post(`${apiUrl}/reboot`);
  //  alert(data.message);
  toast.success(data.message);
  } catch (error) {
   handleApiError(error);
  }
 };

 return (
  <div className='mb-3'>
   <p>Actions</p>
   <div className='action-button'>
    <button className='btn btn-danger form-control' id='shutdown' onClick={handleShutdown}>
     shutdown
    </button>
    <p></p>
    <button className='btn btn-warning form-control' id='reboot' onClick={handleReboot}>
     reboot
    </button>
   </div>
  </div>
 );
};

export default Action;
