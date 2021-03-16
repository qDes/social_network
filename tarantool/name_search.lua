function name_search(first, second)
    local ret = {}
    for _, tuple in box.space.users.index.secondary:pairs({second, first}, {iterator='GE'}) do
        if (string.startswith(tuple[3], first, 1, -1) and string.startswith(tuple[4], second, 1, -1)) then
            table.insert(ret, tuple)
        end
    end
    return ret
end