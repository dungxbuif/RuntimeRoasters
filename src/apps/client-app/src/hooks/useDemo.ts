import { useQuery } from "@tanstack/react-query";
import { demoService } from "@/services/demo.service";

export const useDemoPing = (enabled = false) => {
  return useQuery({
    queryKey: ['demo-ping'],
    queryFn: demoService.ping,
    enabled,
  });
};
