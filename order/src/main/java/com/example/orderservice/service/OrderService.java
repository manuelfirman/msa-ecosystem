package com.example.orderservice.service;
import java.util.Optional;
import java.util.stream.Collectors;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.example.orderservice.repository.OrderRepository;
import com.example.orderservice.model.Order;
import com.example.orderservice.model.OrderItem;
import com.example.orderservice.dto.OrderDTO;
import com.example.orderservice.dto.OrderItemDTO;

import java.math.BigDecimal;
import java.util.List;


@Service
public class OrderService {

     @Autowired
    private OrderRepository orderRepository;

    @Transactional
    public OrderDTO createOrder(OrderDTO orderDTO) {
        Order order = convertToEntity(orderDTO);
        order.setTotalPrice(calculateTotalPrice(order.getOrderItems()));
        for (OrderItem item : order.getOrderItems()) {
            item.setOrder(order);
        }
        Order savedOrder = orderRepository.save(order);
        return convertToDTO(savedOrder);
    }

    @Transactional
    public OrderDTO updateOrder(Long orderId, OrderDTO orderDTO) {
        Optional<Order> existingOrderOpt = orderRepository.findById(orderId);
        if (!existingOrderOpt.isPresent()) {
            throw new RuntimeException("Order not found");
        }

        Order existingOrder = existingOrderOpt.get();
        // Update fields as needed
        existingOrder.setUserId(orderDTO.getUserId());
        existingOrder.setOrderDate(orderDTO.getOrderDate());
        existingOrder.setTotalPrice(orderDTO.getTotalPrice());
        existingOrder.setStatus(orderDTO.getStatus());
        existingOrder.setPaymentMethod(orderDTO.getPaymentMethod());
        existingOrder.setPaymentStatus(orderDTO.getPaymentStatus());
        existingOrder.setOrderItems(orderDTO.getOrderItems().stream()
            .map(this::convertToEntity)
            .collect(Collectors.toList()));

        Order updatedOrder = orderRepository.save(existingOrder);
        return convertToDTO(updatedOrder);
    }

    @Transactional
    public void deleteOrder(Long orderId) {
        orderRepository.deleteById(orderId);
    }

    public OrderDTO getOrderById(Long orderId) {
        Optional<Order> orderOpt = orderRepository.findById(orderId);
        if (!orderOpt.isPresent()) {
            throw new RuntimeException("Order not found");
        }
        return convertToDTO(orderOpt.get());
    }

    public List<OrderDTO> getAllOrders() {
        List<Order> orders = orderRepository.findAll();
        return orders.stream()
            .map(this::convertToDTO)
            .collect(Collectors.toList());
    }

    private BigDecimal calculateTotalPrice(List<OrderItem> orderItems) {
        if (orderItems == null || orderItems.isEmpty()) {
            return BigDecimal.ZERO;
        }
    
        return orderItems.stream()
        .filter(item -> item.getPrice() != null && item.getQuantity() != null)
        .map(item -> item.getPrice().multiply(BigDecimal.valueOf(item.getQuantity())))
        .reduce(BigDecimal.ZERO, BigDecimal::add);
    }

    private Order convertToEntity(OrderDTO orderDTO) {
        Order order = new Order();
        order.setId(orderDTO.getId());
        order.setUserId(orderDTO.getUserId());
        order.setOrderDate(orderDTO.getOrderDate());
        order.setTotalPrice(orderDTO.getTotalPrice());
        order.setStatus(orderDTO.getStatus());
        order.setPaymentMethod(orderDTO.getPaymentMethod());
        order.setPaymentStatus(orderDTO.getPaymentStatus());
    
        List<OrderItem> orderItems = orderDTO.getOrderItems().stream()
            .map(this::convertToEntity)
            .collect(Collectors.toList());
    
        // Establece la relación bidireccional
        for (OrderItem item : orderItems) {
            item.setOrder(order); // Establece la relación en el lado de OrderItem
        }
    
        order.setOrderItems(orderItems);
        return order;
    }
    
    private OrderItem convertToEntity(OrderItemDTO dto) {
        OrderItem orderItem = new OrderItem();
        orderItem.setProductId(dto.getProductId());
        orderItem.setQuantity(dto.getQuantity());
        orderItem.setPrice(dto.getPrice());
        return orderItem;
    }

    private OrderDTO convertToDTO(Order order) {
        OrderDTO orderDTO = new OrderDTO();
        orderDTO.setId(order.getId());
        orderDTO.setUserId(order.getUserId());
        orderDTO.setOrderDate(order.getOrderDate());
        orderDTO.setTotalPrice(order.getTotalPrice());
        orderDTO.setStatus(order.getStatus());
        orderDTO.setPaymentMethod(order.getPaymentMethod());
        orderDTO.setPaymentStatus(order.getPaymentStatus());
        orderDTO.setOrderItems(order.getOrderItems().stream()
            .map(this::convertToDTO)
            .collect(Collectors.toList()));
        return orderDTO;
    }

    private OrderItemDTO convertToDTO(OrderItem orderItem) {
        OrderItemDTO dto = new OrderItemDTO();
        dto.setProductId(orderItem.getProductId());
        dto.setQuantity(orderItem.getQuantity());
        dto.setPrice(orderItem.getPrice());
        return dto;
    }

}
