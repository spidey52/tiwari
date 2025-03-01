import { AxiosError } from "axios";
import { toast } from "react-toastify";

const handleApiError = (error: unknown) => {
 if (error instanceof AxiosError) {
  if (error.response?.data) {
   //  alert(JSON.stringify(error.response.data));
   toast.error(JSON.stringify(error.response.data));
  } else {
   //  alert(error.message);
   toast.error(error.message);
  }

  return;
 }

 if (error instanceof Error) {
  console.error(error.message);
 }
};

export default handleApiError;
