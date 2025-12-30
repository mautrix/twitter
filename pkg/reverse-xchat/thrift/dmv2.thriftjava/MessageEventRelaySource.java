package com.x.dmv2.thriftjava;

import kotlin.Metadata;
import kotlin.enums.EnumEntries;
import kotlin.enums.EnumEntriesKt;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

/* JADX WARN: Failed to restore enum class, 'enum' modifier and super class removed */
/* JADX WARN: Unknown enum class pattern. Please report as an issue! */
@Metadata(m64929d1 = {"\u0000\u0012\n\u0002\u0018\u0002\n\u0002\u0010\u0010\n\u0000\n\u0002\u0010\b\n\u0002\b\u0006\b\u0086\u0081\u0002\u0018\u0000 \b2\b\u0012\u0004\u0012\u00020\u00000\u0001:\u0001\bB\u0011\b\u0002\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005R\u0010\u0010\u0002\u001a\u00020\u00038\u0006X\u0087\u0004¢\u0006\u0002\n\u0000j\u0002\b\u0006j\u0002\b\u0007¨\u0006\t"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventRelaySource;", "", "value", "", "<init>", "(Ljava/lang/String;II)V", "LiveFanout", "MessagePull", "Companion", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final class MessageEventRelaySource {
    private static final /* synthetic */ EnumEntries $ENTRIES;
    private static final /* synthetic */ MessageEventRelaySource[] $VALUES;

    /* renamed from: Companion, reason: from kotlin metadata */
    @InterfaceC88464a
    public static final Companion INSTANCE;
    public static final MessageEventRelaySource LiveFanout = new MessageEventRelaySource("LiveFanout", 0, 0);
    public static final MessageEventRelaySource MessagePull = new MessageEventRelaySource("MessagePull", 1, 1);

    @JvmField
    public final int value;

    @Metadata(m64929d1 = {"\u0000\u0018\n\u0002\u0018\u0002\n\u0002\u0010\u0000\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\u0003\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0010\u0010\u0004\u001a\u0004\u0018\u00010\u00052\u0006\u0010\u0006\u001a\u00020\u0007¨\u0006\b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventRelaySource$Companion;", "", "<init>", "()V", "findByValue", "Lcom/x/dmv2/thriftjava/MessageEventRelaySource;", "value", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class Companion {
        public /* synthetic */ Companion(DefaultConstructorMarker defaultConstructorMarker) {
            this();
        }

        @InterfaceC88465b
        public final MessageEventRelaySource findByValue(int value) {
            if (value == 0) {
                return MessageEventRelaySource.LiveFanout;
            }
            if (value != 1) {
                return null;
            }
            return MessageEventRelaySource.MessagePull;
        }

        private Companion() {
        }
    }

    private static final /* synthetic */ MessageEventRelaySource[] $values() {
        return new MessageEventRelaySource[]{LiveFanout, MessagePull};
    }

    static {
        MessageEventRelaySource[] messageEventRelaySourceArr$values = $values();
        $VALUES = messageEventRelaySourceArr$values;
        $ENTRIES = EnumEntriesKt.m65223a(messageEventRelaySourceArr$values);
        INSTANCE = new Companion(null);
    }

    private MessageEventRelaySource(String str, int i, int i2) {
        this.value = i2;
    }

    @InterfaceC88464a
    public static EnumEntries getEntries() {
        return $ENTRIES;
    }

    public static MessageEventRelaySource valueOf(String str) {
        return (MessageEventRelaySource) Enum.valueOf(MessageEventRelaySource.class, str);
    }

    public static MessageEventRelaySource[] values() {
        return (MessageEventRelaySource[]) $VALUES.clone();
    }
}
