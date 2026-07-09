import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchMyProvider as fetchMyProviderApi,
  fetchProvider as fetchProviderByIdApi,
} from "../api/providers";
import {
  fetchProviderServices as fetchProviderServicesApi,
  createService as createServiceApi,
  updateService as updateServiceApi,
  updateServiceStatus as updateServiceStatusApi,
  replaceServiceTags,
  type CreateServiceInput,
  type UpdateServiceInput,
  type ServiceQueryParams,
} from "../api/services";
import { useAuthStore } from "../stores/auth.store";

export function useMyProvider() {
  return useQuery({
    queryKey: ["my-provider"],
    queryFn: fetchMyProviderApi,
    retry: false,
  });
}

export function useProviderById(id: number | undefined) {
  return useQuery({
    queryKey: ["provider", id],
    queryFn: () => fetchProviderByIdApi(id!),
    enabled: !!id,
  });
}

export function useProviderServices(
  providerId: number | undefined,
  params: ServiceQueryParams = {}
) {
  return useQuery({
    queryKey: ["provider-services", providerId, params],
    queryFn: () => fetchProviderServicesApi(providerId!, params),
    enabled: !!providerId,
  });
}

export function useCreateService() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateServiceInput) => createServiceApi(data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["provider-services"] });
    },
  });
}

export function useUpdateService() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateServiceInput }) =>
      updateServiceApi(id, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["provider-services"] });
    },
  });
}

export function useReplaceServiceTags() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      tagIds,
    }: {
      id: number;
      tagIds: number[];
    }) => replaceServiceTags(id, tagIds),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["provider-services"] });
    },
  });
}

export function useUpdateServiceStatus() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      status,
    }: {
      id: number;
      status: "active" | "inactive";
    }) => updateServiceStatusApi(id, status),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["provider-services"] });
    },
  });
}
